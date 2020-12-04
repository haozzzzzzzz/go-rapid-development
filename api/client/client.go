package client

import (
	"context"

	"fmt"

	http2 "net/http"

	"github.com/haozzzzzzzz/go-rapid-development/api/code"
	"github.com/haozzzzzzzz/go-rapid-development/api/request"
	"github.com/haozzzzzzzz/go-rapid-development/http"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Ctx        context.Context
	HttpClient *http2.Client
	UrlPrefix  string
}

func NewClient(
	ctx context.Context,
	httpClient *http2.Client,
	urlPrefix string,
) *Client {
	return &Client{
		Ctx:        ctx,
		HttpClient: httpClient,
		UrlPrefix:  urlPrefix,
	}
}

func (m *Client) Get(
	urlPath string,
	iRespData interface{},
	iPathData interface{},
	iQueryData interface{},
) (err error) {
	defer func() {
		iRecover := recover()
		if iRecover != nil {
			err = uerrors.Newf("panic: %#v", iRecover)
		}
	}()

	apiUrl := fmt.Sprintf("%s%s", m.UrlPrefix, urlPath)
	reqUrl, err := http.NewUrlByStrUrl(apiUrl)
	if nil != err {
		logrus.Errorf("new url by string url failed. %s.", err)
		return
	}

	pathData := make(map[string][]string)
	err = request.ObjToUrlValues(pathData, iPathData)
	if nil != err {
		logrus.Errorf("map path data failed. %s.", err)
		return
	}

	queryData := make(map[string][]string)
	err = request.ObjToUrlValues(queryData, iQueryData)
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

func (m *Client) GetWithHeader(
	urlPath string,
	iRespData interface{},
	iPathData interface{},
	iQueryData interface{},
	header http2.Header,
) (err error) {
	defer func() {
		iRecover := recover()
		if iRecover != nil {
			err = uerrors.Newf("panic: %#v", iRecover)
		}
	}()

	apiUrl := fmt.Sprintf("%s%s", m.UrlPrefix, urlPath)
	reqUrl, err := http.NewUrlByStrUrl(apiUrl)
	if nil != err {
		logrus.Errorf("new url by string url failed. %s.", err)
		return
	}

	pathData := make(map[string][]string)
	err = request.ObjToUrlValues(pathData, iPathData)
	if nil != err {
		logrus.Errorf("map path data failed. %s.", err)
		return
	}

	queryData := make(map[string][]string)
	err = request.ObjToUrlValues(queryData, iQueryData)
	if nil != err {
		logrus.Errorf("map query data failed. %s.", err)
		return
	}

	reqUrl.SetPathData(pathData)
	reqUrl.SetQueryData(queryData)

	req := http.NewRequestByUrl(reqUrl, m.Ctx, m.HttpClient)
	req.Header = header
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
	defer func() {
		iRecover := recover()
		if iRecover != nil {
			err = uerrors.Newf("panic: %#v", iRecover)
		}
	}()

	apiUrl := fmt.Sprintf("%s%s", m.UrlPrefix, urlPath)
	reqUrl, err := http.NewUrlByStrUrl(apiUrl)
	if nil != err {
		logrus.Errorf("new url by string url failed. %s.", err)
		return
	}

	pathData := make(map[string][]string)
	err = request.ObjToUrlValues(pathData, iPathData)
	if nil != err {
		logrus.Errorf("map path data failed. %s.", err)
		return
	}

	queryData := make(map[string][]string)
	err = request.ObjToUrlValues(queryData, iQueryData)
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

func (m *Client) PostWithHeader(
	urlPath string,
	iRespData interface{},
	iPathData interface{},
	iQueryData interface{},
	iPostData interface{},
	header http2.Header,
) (err error) {
	defer func() {
		iRecover := recover()
		if iRecover != nil {
			err = uerrors.Newf("panic: %#v", iRecover)
		}
	}()

	apiUrl := fmt.Sprintf("%s%s", m.UrlPrefix, urlPath)
	reqUrl, err := http.NewUrlByStrUrl(apiUrl)
	if nil != err {
		logrus.Errorf("new url by string url failed. %s.", err)
		return
	}

	pathData := make(map[string][]string)
	err = request.ObjToUrlValues(pathData, iPathData)
	if nil != err {
		logrus.Errorf("map path data failed. %s.", err)
		return
	}

	queryData := make(map[string][]string)
	err = request.ObjToUrlValues(queryData, iQueryData)
	if nil != err {
		logrus.Errorf("map query data failed. %s.", err)
		return
	}

	reqUrl.SetPathData(pathData)
	reqUrl.SetQueryData(queryData)

	req := http.NewRequestByUrl(reqUrl, m.Ctx, m.HttpClient)
	req.Header = header
	err = req.PostJson(iPostData, iRespData)
	if nil != err {
		logrus.Errorf("request api json failed. %s.", err)
		return
	}

	return
}

func (m *Client) PostJsonReq(
	urlPath string,
	iRespData interface{},
	iPathData interface{},
	iQueryData interface{},
	iPostData interface{},
	beforeRequest func(rawReq *http2.Request) (err error),
) (err error) {
	apiUrl := fmt.Sprintf("%s%s", m.UrlPrefix, urlPath)
	reqUrl, err := http.NewUrlByStrUrl(apiUrl)
	if nil != err {
		logrus.Errorf("new url by string url failed. %s.", err)
		return
	}

	pathData := make(map[string][]string)
	err = request.ObjToUrlValues(pathData, iPathData)
	if nil != err {
		logrus.Errorf("map path data failed. %s.", err)
		return
	}

	queryData := make(map[string][]string)
	err = request.ObjToUrlValues(queryData, iQueryData)
	if nil != err {
		logrus.Errorf("map query data failed. %s.", err)
		return
	}

	reqUrl.SetPathData(pathData)
	reqUrl.SetQueryData(queryData)

	req := http.NewRequestByUrl(reqUrl, m.Ctx, m.HttpClient)
	rawReq, err := req.BuildPostJson(iPostData)
	if err != nil {
		logrus.Errorf("req build post json failed. error: %s", err)
		return
	}

	if beforeRequest != nil {
		err = beforeRequest(rawReq)
		if err != nil {
			logrus.Errorf("call before request failed. error: %s", err)
			return
		}
	}

	err = req.DoRespJson(rawReq, iRespData)
	if err != nil {
		logrus.Errorf("req do req json failed. error: %s", err)
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
