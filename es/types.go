package es

import (
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
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
