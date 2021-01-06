package proxy

import (
	"errors"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uctx"
	"github.com/sirupsen/logrus"
	"github.com/wangjia184/sortedset"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// 等待分配代理url超时
var ErrorWaitProxyUrlTimeout = errors.New("wait proxy url timeout")

const DefaultNoProxyReqHeaderKey = "X-Do-No-Use-Proxy"

func RequestDoNotUserProxy(req *http.Request) {
	req.Header.Set(DefaultNoProxyReqHeaderKey, "true")
}

// proxy url dispatcher
// 1. dispatch proxy url with per url rate limit
// 2. support proxy ignore url list
// 3. support request not use proxy
// 3. priority request queue

type Dispatcher struct {
	NotProxyReqHeaderKey string   // request header specify do not use proxy. value: 1/true: yes
	NotProxyReqUrls      []string // request url specify do not use proxy

	ProxyUrls *ProxyUrls
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		NotProxyReqHeaderKey: DefaultNoProxyReqHeaderKey,
		NotProxyReqUrls:      []string{},
		ProxyUrls:            NewProxyUrls(),
	}
}

// allot proxy url for request
func (m *Dispatcher) DispatchProxyUrl(
	req *http.Request,
	waitTimeout time.Duration, // <0: not timeout; =0: not wait; >0: wait
) (
	success bool,
	urlUrl *url.URL,
	err error,
) {
	if m.DoNotUserProxy(req) {
		return
	}

	success, strUrl, err := m.ProxyUrls.SelectUrl(waitTimeout)
	if err != nil {
		logrus.Errorf("select proxy url failed. error: %s", err)
		return
	}

	urlUrl, err = url.ParseRequestURI(strUrl)
	if err != nil {
		success = false
		logrus.Errorf("parse proxy url failed. error: %s", err)
		return
	}

	return
}

func (m *Dispatcher) DoNotUserProxy(req *http.Request) (notUseProxy bool) {
	notProxyHeaderVal := req.Header.Get(m.NotProxyReqHeaderKey)
	notProxyHeaderVal = strings.ToLower(notProxyHeaderVal)
	if notProxyHeaderVal == "1" || notProxyHeaderVal == "true" {
		notUseProxy = true
		return
	}

	if len(m.NotProxyReqUrls) <= 0 {
		return
	}

	reqUrl := req.URL.String()
	for _, notProxyUrl := range m.NotProxyReqUrls {
		if strings.Contains(reqUrl, notProxyUrl) {
			notUseProxy = true
			return
		}
	}

	return
}

type ProxyUrl struct {
	Url     string
	Limiter *rate.Limiter // rate limiter
}

func NewProxyUrl(url string, limiter *rate.Limiter) *ProxyUrl {
	return &ProxyUrl{
		Url:     url,
		Limiter: limiter,
	}
}

// proxy urls sorted set
type ProxyUrls struct {
	sortedSet *sortedset.SortedSet
	lock      sync.RWMutex
}

func NewProxyUrls() *ProxyUrls {
	return &ProxyUrls{
		sortedSet: sortedset.New(),
	}
}

func (m *ProxyUrls) AddUrl(proxyUrl *ProxyUrl) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.sortedSet.AddOrUpdate(proxyUrl.Url, 1, proxyUrl)
}

func (m *ProxyUrls) RemoveUrl(urlKey string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.sortedSet.Remove(urlKey)
}

func (m *ProxyUrls) chooseMinScoreUrl() (proxyUrl *ProxyUrl) {
	m.lock.Lock()
	defer m.lock.Unlock()
	node := m.sortedSet.GetByRank(1, false)
	if node == nil {
		return
	}

	nodeKey := node.Key()
	proxyUrl, ok := node.Value.(*ProxyUrl)
	if !ok {
		m.sortedSet.Remove(nodeKey)
		return
	}

	m.sortedSet.AddOrUpdate(nodeKey, node.Score()+1, node.Value)

	return
}

func (m *ProxyUrls) SelectUrl(
	waitTimeout time.Duration, // <0: not timeout; =0: not wait; >0: wait
) (
	success bool,
	strUrl string,
	err error,
) {

	proxyUrl := m.chooseMinScoreUrl()
	if proxyUrl == nil || proxyUrl.Url == "" {
		return
	}

	if proxyUrl.Limiter != nil {
		if waitTimeout > 0 {
			ctx, cancel := uctx.BackgroundWithTimeout(waitTimeout)
			defer cancel()
			err = proxyUrl.Limiter.Wait(ctx)
			if err != nil {
				logrus.Errorf("proxy url limiter wait failed. url: %s, error: %s", proxyUrl.Url, err)
				err = ErrorWaitProxyUrlTimeout
				return
			}

		} else if waitTimeout < 0 {
			ctx, cancel := uctx.BackgroundWithCancel()
			defer cancel()
			err = proxyUrl.Limiter.Wait(ctx)
			if err != nil {
				logrus.Errorf("proxy url limiter wait failed. url: %s, error: %s", proxyUrl.Url, err)
				return
			}

		} else {
			allow := proxyUrl.Limiter.Allow()
			if !allow {
				err = errors.New("proxy deny")
				return
			}

		}

	}

	strUrl = proxyUrl.Url
	success = true
	return
}
