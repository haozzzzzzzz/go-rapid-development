package es

import (
	"encoding/json"
	"github.com/haozzzzzzzz/go-rapid-development/v2/utils/uerrors"
	"github.com/sirupsen/logrus"
)

type ErrorResponseBody struct {
	Error struct {
		Type   string `json:"type"`
		Reason string `json:"reason"`
	} `json:"error"`
}

func (m *ErrorResponseBody) Err() error {
	return uerrors.Newf("response error, %s:%s", m.Error.Type, m.Error.Reason)
}

func ParseErrorResponseBody(
	bBody []byte,
) (body *ErrorResponseBody, err error) {
	body = &ErrorResponseBody{}
	if bBody == nil {
		return
	}

	err = json.Unmarshal(bBody, body)
	if nil != err {
		logrus.Errorf("json unmarshal reponse body failed. error: %s.", err)
		return
	}
	return
}

const ErrorTypeDocumentMissing = "document_missing_exception" // 文档不存在
