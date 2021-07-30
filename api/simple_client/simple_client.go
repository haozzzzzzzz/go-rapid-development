package simple_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/haozzzzzzzz/go-rapid-development/v2/utils/umapstructure"
	"github.com/sirupsen/logrus"
	"github.com/xiam/to"
	"golang.org/x/net/context/ctxhttp"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	ShortTimeoutRequestClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     5 * time.Minute,
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

type ApiClient struct {
	UrlPrefix  string
	Ctx        context.Context
	HttpClient *http.Client
}

type ReqParams struct {
	Query  interface{} // url.Values、map[string]interface{}、struct with json tag
	Body   interface{} // json.Marshal
	Header interface{} // http.Header
}

func (m *ApiClient) DoRequest(
	method string,
	uri string,
	params *ReqParams,
	respBody interface{}, // json.Unmarshal
) (err error) {
	strUrl := fmt.Sprintf("%s%s", m.UrlPrefix, uri)
	logger := logrus.WithFields(logrus.Fields{
		"func":   "ApiClient.DoRequest",
		"url":    strUrl,
		"params": params,
	})

	var bodyReader io.Reader
	reqBody := params.Body
	query := params.Query
	header, err := parseHeader(params.Header)
	if err != nil {
		logger.Errorf("parse header failed. %s", err)
		return
	}

	if header == nil {
		header = make(http.Header)
	}

	if reqBody != nil {
		bBody, errM := json.Marshal(reqBody)
		err = errM
		if err != nil {
			logger.Errorf("marshal json reqBody failed. %s", err)
			return
		}
		bodyReader = bytes.NewReader(bBody)
		header.Set("Content-Type", "application/json")
	}

	oUrl, err := url.Parse(strUrl)
	if err != nil {
		logger.Errorf("parse str url failed. %s", err)
		return
	}

	err = parseQuery(oUrl, query)
	if err != nil {
		logger.Errorf("parse query failed. error: %s", err)
		return
	}

	req, err := http.NewRequest(method, oUrl.String(), bodyReader)
	if err != nil {
		return
	}

	for key, values := range header {
		req.Header[key] = values
	}

	resp, err := ctxhttp.Do(m.Ctx, m.HttpClient, req)
	if err != nil {
		logger.Errorf("do request failed. %s", err)
		return
	}

	defer func() {
		errC := resp.Body.Close()
		if errC != nil {
			logger.Errorf("close resp body failed.")
		}
	}()

	bRespBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("read all resp body failed. %s", err)
		return
	}

	err = json.Unmarshal(bRespBody, respBody)
	if err != nil {
		logger.Errorf("unmarshal resp body failed. err: %s body: %s", err, string(bRespBody))
		return
	}

	return
}

type ApiCode struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg"`
}

func (m *ApiCode) Error() (err error) {
	if m.Ret != ApiCodeSuccess {
		err = fmt.Errorf("api return error. ret: %d, msg: %s", m.Ret, m.Msg)
		return
	}
	return
}

const ApiCodeSuccess int = 1

// 解析query数据到url
// 支持url.Values, map[string]interface, 带json tag的struct结构体
func parseQuery(
	oUrl *url.URL,
	query interface{},
) (
	err error,
) {
	if query == nil {
		return
	}

	existQuery := oUrl.Query()
	queryMap := make(map[string]interface{}, 0)
	switch t := query.(type) {
	case url.Values:
		for key, values := range t {
			existQuery[key] = values
		}

	case map[string]interface{}:
		queryMap = t

	default:
		err = umapstructure.Decode(t, &queryMap, "json")
		if err != nil {
			logrus.Errorf("decode mapstructure failed. error: %s", err)
			return
		}
	}

	for key, value := range queryMap {
		existQuery[key] = []string{
			to.String(value),
		}
	}

	oUrl.RawQuery = existQuery.Encode()

	return
}

func parseHeader(iHeader interface{}) (header http.Header, err error) {
	logger := logrus.WithFields(logrus.Fields{
		"func":   "parseHeader",
		"header": iHeader,
	})
	header = make(http.Header)
	headerMap := make(map[string]interface{})
	switch t := iHeader.(type) {
	case http.Header:
		header = t
		return
	case map[string]interface{}:
		headerMap = t
	default:
		err = umapstructure.Decode(t, &headerMap, "json")
		if err != nil {
			logger.Errorf("decode mapstructure failed. error: %s", err)
			return
		}
	}

	for key, value := range headerMap {
		header.Set(key, to.String(value))
	}
	return
}
