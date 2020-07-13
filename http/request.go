package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/haozzzzzzzz/go-rapid-development/utils/ujson"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
)

var defaultRequestCheckerMaker RequestCheckerMaker

func SetDefaultRequestCheckerMaker(checkerMaker RequestCheckerMaker) {
	defaultRequestCheckerMaker = checkerMaker
}

type RequestChecker interface {
	Before(request *Request, rawReq *http.Request)
	After(response *http.Response, err error)
}

type RequestCheckerMaker interface {
	NewChecker() RequestChecker
}

type Request struct {
	Url                 *Url // 必填
	Header              http.Header
	codecs              string
	Ctx                 context.Context
	Client              *http.Client
	RequestCheckerMaker RequestCheckerMaker
}

func (m *Request) RequestChecker() RequestChecker {
	if m.RequestCheckerMaker == nil {
		return nil
	}

	return m.RequestCheckerMaker.NewChecker()
}

func NewRequest(
	strUrl string,
	ctx context.Context,
	client *http.Client,
) (req *Request, err error) {
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
		Url:                 url,
		Header:              make(http.Header),
		codecs:              "json",
		Ctx:                 ctx,
		Client:              client,
		RequestCheckerMaker: defaultRequestCheckerMaker,
	}

	return
}

func NewRequestByUrl(
	reqUrl *Url,
	ctx context.Context,
	client *http.Client,
) (req *Request) {
	if client == nil {
		client = RequestClient
	}

	req = &Request{
		Url:                 reqUrl,
		Header:              make(http.Header),
		codecs:              "json",
		Ctx:                 ctx,
		Client:              client,
		RequestCheckerMaker: defaultRequestCheckerMaker,
	}

	return
}

func (m *Request) URL() string {
	return m.Url.String()
}

// when err !=nil, resp returns nil
func (m *Request) Get() (resp *http.Response, err error) {
	strUrl := m.URL()

	req, err := http.NewRequest("GET", strUrl, nil)
	if err != nil {
		logrus.Errorf("new http request failed. error: %s.", err)
		return
	}

	req.Header = m.Header
	host := m.Header.Get("Host")
	if host != "" {
		req.Host = host
	}

	resp, err = m.Do(req)
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

	err = ujson.UnmarshalJsonFromReader(ack.Body, v)
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
	var response *http.Response

	var bytesBody []byte
	if body != nil {
		bytesBody, err = json.Marshal(body)
		if err != nil {
			logrus.Warnf("marshal post body failed. %s", err)
			return
		}
	}

	strUrl := m.URL()
	bodyReader := bytes.NewBuffer(bytesBody)
	contentType := "application/json;charset=utf8"

	req, err := http.NewRequest("POST", strUrl, bodyReader)
	if err != nil {
		return
	}
	req.Header = m.Header
	req.Header.Set("Content-Type", contentType)

	response, err = m.Do(req)
	if err != nil {
		logrus.Warnf("post request failed. %s", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("http status: %s", response.Status))
		logrus.Errorf("response error. %s.", err)
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

	err = ujson.UnmarshalJsonFromReader(response.Body, resp)
	if nil != err {
		logrus.Errorf("unmarshal body json failed. %s.", err)
		return
	}

	return
}

func (m *Request) Do(rawRequest *http.Request) (response *http.Response, err error) {
	// checker
	checker := m.RequestChecker()
	if checker != nil {
		checker.Before(m, rawRequest)
		defer func() {
			checker.After(response, err)
		}()
	}

	if m.Ctx != nil {
		response, err = ctxhttp.Do(m.Ctx, m.Client, rawRequest)
	} else {
		response, err = m.Client.Do(rawRequest)
	}

	if nil != err {
		logrus.Errorf("do request failed. error: %s.", err)
		return
	}

	return
}
