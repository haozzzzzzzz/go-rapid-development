package uerrors

import (
	"errors"
	"fmt"
)

func Newf(format string, values ...interface{}) error {
	return errors.New(fmt.Sprintf(format, values...))
}

type StringError string

func (m StringError) Error() string {
	return string(m)
}

// 查询不到记录
const ErrQueryNoRecord StringError = "Err query no record"
