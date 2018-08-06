package client

import (
	"context"

	"fmt"

	http2 "net/http"

	"github.com/haozzzzzzzz/go-rapid-development/api/code"
	"github.com/haozzzzzzzz/go-rapid-development/api/request"
	"github.com/haozzzzzzzz/go-rapid-development/utils/http"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Ctx        context.Context
	HttpClient *http2.Client
	UrlPrefix  string
}

func (m *Client) Get(
	urlPath string,
	iRespData interface{},
	iPathData interface{},
	iQueryData interface{},
) (err error) {
	apiUrl := fmt.Sprintf("%s%s", m.UrlPrefix, urlPath)
	reqUrl, err := http.NewUrlByStrUrl(apiUrl)
	if nil != err {
		logrus.Errorf("new url by string url failed. %s.", err)
		return
	}

	pathData := make(map[string][]string)
	err = request.FormMapStruct(pathData, iPathData)
	if nil != err {
		logrus.Errorf("map path data failed. %s.", err)
		return
	}

	queryData := make(map[string][]string)
	err = request.FormMapStruct(queryData, iQueryData)
	if nil != err {
		logrus.Errorf("map query data failed. %s.", err)
		return
	}

	reqUrl.SetPathData(pathData)
	reqUrl.SetQueryData(queryData)

	req := http.NewRequestByUrl(reqUrl, m.Ctx, m.HttpClient)
	err = req.GetJSON(iRespData)
	if nil != err {
		logrus.Errorf("request api json failed. %s.", err)
		return
	}

	return
}

func (m *Client) Post(
	urlPath string,
	iRespData interface{},
	iPathData interface{},
	iQueryData interface{},
	iPostData interface{},
) (err error) {
	apiUrl := fmt.Sprintf("%s%s", m.UrlPrefix, urlPath)
	reqUrl, err := http.NewUrlByStrUrl(apiUrl)
	if nil != err {
		logrus.Errorf("new url by string url failed. %s.", err)
		return
	}

	pathData := make(map[string][]string)
	err = request.FormMapStruct(pathData, iPathData)
	if nil != err {
		logrus.Errorf("map path data failed. %s.", err)
		return
	}

	queryData := make(map[string][]string)
	err = request.FormMapStruct(queryData, iQueryData)
	if nil != err {
		logrus.Errorf("map query data failed. %s.", err)
		return
	}

	reqUrl.SetPathData(pathData)
	reqUrl.SetQueryData(queryData)

	req := http.NewRequestByUrl(reqUrl, m.Ctx, m.HttpClient)
	err = req.PostJson(iPostData, iRespData)
	if nil != err {
		logrus.Errorf("request api json failed. %s.", err)
		return
	}

	return
}

type ResponseMessage struct {
	ReturnCode uint32 `json:"ret"`
	Message    string `json:"msg"`
}

func (m *ResponseMessage) ApiCode() *code.ApiCode {
	return &code.ApiCode{
		Code:    m.ReturnCode,
		Message: m.Message,
	}
}
