package es

import (
	"errors"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uio"
	"github.com/sirupsen/logrus"
)

func RespError(
	errRespBody []byte,
) (respBody *ErrorResponseBody, err error) {
	respBody, err = ParseErrorResponseBody(errRespBody)
	if nil != err {
		logrus.Errorf("parse error response body failed. error: %s.", err)
		return
	}

	if respBody == nil {
		err = errors.New("parse error response body get nil")
		return
	}

	err = respBody.Err()

	return
}

// handle response body
func ReadRespError(response *esapi.Response) (errBody *ErrorResponseBody, err error) {
	bBody, err := uio.ReadAllAndClose(response.Body)
	if err != nil {
		logrus.Errorf("read response body failed. error: %s", err)
		return
	}

	errBody, err = ParseErrorResponseBody(bBody)
	if nil != err {
		logrus.Errorf("parse error response body failed. error: %s.", err)
		return
	}

	if errBody == nil {
		err = errors.New("read resp error body but got nil")
		return
	}
	return
}
