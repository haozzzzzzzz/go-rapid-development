package es

import (
	"errors"
	"github.com/sirupsen/logrus"
)

func RespError(errRespBody []byte) (err error) {
	respBody, err := ParseErrorResponseBody(errRespBody)
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
