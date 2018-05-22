package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	json2 "github.com/haozzzzzzzz/go-rapid-development/utils/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
)

var (
	ClientRequest = &http.Client{
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
	url    string
	query  map[string]string
	codecs string
	Ctx    context.Context
	Client *http.Client
}

func NewRequest(url string, ctx context.Context, client *http.Client) *Request {
	if client == nil {
		client = ClientRequest
	}

	return &Request{
		url:    url,
		codecs: "json",
		query:  make(map[string]string),
		Ctx:    ctx,
		Client: client,
	}
}

func (m *Request) URL() string {

	if len(m.query) == 0 {
		return m.url
	} else {
		values := make([]string, 0)

		for key, value := range m.query {
			values = append(values, fmt.Sprintf("%s=%s", key, value))
		}

		return fmt.Sprintf(
			"%s?%s",
			m.url,
			strings.Join(values, "&"),
		)
	}
}

func (m *Request) Query(name string, value interface{}) {
	m.query[name] = fmt.Sprintf("%v", value)
}

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
		logrus.Errorf("response error. %s.", err)
		return
	}

	return
}

func (m *Request) GetJSON(v interface{}) (err error) {
	ack, err := m.Get()
	if err != nil {
		return
	}

	err = json2.UnmarshalJsonFromReader(ack.Body, v)
	return
}

func (m *Request) GetText() (text string, err error) {
	resp, err := m.Get()
	if err != nil {
		return
	}

	bytesBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	text = string(bytesBody)
	return
}

func (m *Request) PostJson(body interface{}) (respBody []byte, err error) {
	bytesBody, err := json.Marshal(body)
	if err != nil {
		logrus.Warnf("marshal post body failed. %s", err)
		return
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

	respBody, err = ioutil.ReadAll(response.Body)
	if err != nil {
		logrus.Warnf("read response body failed. %s", err)
		return
	}

	return
}
