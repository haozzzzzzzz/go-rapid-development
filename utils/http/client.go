package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"fmt"

	json2 "github.com/haozzzzzzzz/go-rapid-development/utils/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
)

var (
	RequestClient = &http.Client{
		Transport: &http.Transport{
			//MaxIdleConns:        1000,
			//MaxIdleConnsPerHost: 200,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				d := net.Dialer{
					Timeout:   2 * time.Second,
					KeepAlive: 10 * time.Minute,
				}

				return d.DialContext(ctx, network, addr)
			},
		},

		Timeout: 2 * time.Second,
	}
)

type Request struct {
	Url    *Url
	codecs string
	Ctx    context.Context
	Client *http.Client
}

func NewRequest(strUrl string, ctx context.Context, client *http.Client) (req *Request, err error) {
	if client == nil {
		client = RequestClient
	}

	url, err := NewUrlByStrUrl(strUrl)
	if nil != err {
		logrus.Errorf("new url failed. %s.", err)
		return
	}

	if client == nil {
		client = RequestClient
	}

	req = &Request{
		Url:    url,
		codecs: "json",
		Ctx:    ctx,
		Client: client,
	}

	return
}

func NewRequestByUrl(reqUrl *Url, ctx context.Context, client *http.Client) (req *Request) {
	if client == nil {
		client = RequestClient
	}

	req = &Request{
		Url:    reqUrl,
		codecs: "json",
		Ctx:    ctx,
		Client: client,
	}
	return
}

func (m *Request) URL() string {
	return m.Url.String()
}

// when err !=nil, resp returns nil
func (m *Request) Get() (resp *http.Response, err error) {
	strUrl := m.URL()
	if m.Ctx != nil {
		resp, err = ctxhttp.Get(m.Ctx, m.Client, strUrl)
	} else {
		resp, err = m.Client.Get(strUrl)
	}

	if err != nil {
		logrus.Errorf("request get failed. %s.", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("http status: %s", resp.Status))
		logrus.Errorf("response error. %s.", err)
		return
	}

	return
}

func (m *Request) GetJSON(v interface{}) (err error) {
	ack, err := m.Get()
	if err != nil {
		logrus.Errorf("request failed. %s.", err)
		return
	}

	defer func() {
		errClose := ack.Body.Close()
		if errClose != nil {
			logrus.Errorf("close http response body failed. %s.", err)
			if err == nil {
				err = errClose
			}
		}
	}()

	err = json2.UnmarshalJsonFromReader(ack.Body, v)
	if nil != err {
		logrus.Errorf("unmarshal body json failed. %s.", err)
		return
	}
	return
}

func (m *Request) GetText() (text string, err error) {
	resp, err := m.Get()
	if err != nil {
		logrus.Errorf("get api failed. %s.", err)
		return
	}

	defer func() {
		errClose := resp.Body.Close()
		if errClose != nil {
			logrus.Errorf("close http response body failed. %s.", err)
			if err == nil {
				err = errClose
			}
		}
	}()

	bytesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("read body failed. %s.", err)
		return
	}
	text = string(bytesBody)
	return
}

func (m *Request) PostJson(body interface{}, resp interface{}) (err error) {
	var bytesBody []byte
	if body != nil {
		bytesBody, err = json.Marshal(body)
		if err != nil {
			logrus.Warnf("marshal post body failed. %s", err)
			return
		}
	}

	strUrl := m.URL()
	contentType := "application/json;charset=utf8"
	var response *http.Response
	if m.Ctx != nil {
		response, err = ctxhttp.Post(m.Ctx, m.Client, strUrl, contentType, bytes.NewBuffer(bytesBody))

	} else {
		response, err = m.Client.Post(strUrl, contentType, bytes.NewBuffer(bytesBody))
	}

	if err != nil {
		logrus.Warnf("post request failed. %s", err)
		return
	}

	defer func() {
		errClose := response.Body.Close()
		if errClose != nil {
			logrus.Errorf("close http response body failed. %s.", err)
			if err == nil {
				err = errClose
			}
		}
	}()

	err = json2.UnmarshalJsonFromReader(response.Body, resp)
	if nil != err {
		logrus.Errorf("unmarshal body json failed. %s.", err)
		return
	}

	return
}

func (m *Request) Do(rawRequest *http.Request) (*http.Response, error) {
	return ctxhttp.Do(m.Ctx, m.Client, rawRequest)
}
