package http

import (
	"net/url"

	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type Url struct {
	url.URL
	QueryValues url.Values
}

func NewUrl(rawUrl *url.URL) *Url {
	return &Url{
		URL:         *rawUrl,
		QueryValues: rawUrl.Query(),
	}
}

func NewUrlByStrUrl(strUrl string) (retUrl *Url, err error) {
	rawUrl, err := url.Parse(strUrl)
	if nil != err {
		logrus.Errorf("parse string url failed. %s.", err)
		return
	}

	retUrl = NewUrl(rawUrl)

	return
}

func (m *Url) GetRawQuery() string {
	return m.QueryValues.Encode()
}

func (m *Url) String() string {
	m.URL.RawQuery = m.GetRawQuery()
	return m.URL.String()
}

// 设置pathData
func (m *Url) SetPathData(pathData url.Values) {
	for key, values := range pathData {
		if len(values) == 0 {
			continue
		}
		m.Path = strings.Replace(m.Path, fmt.Sprintf(":%s", key), values[0], -1)
		m.RawPath = m.EscapedPath()
	}
	return
}

// 设置queryData
func (m *Url) SetQueryData(queryData url.Values) {
	m.QueryValues = queryData
	return
}
