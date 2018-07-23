package uerrors

import (
	"errors"
	"fmt"
)

func Newf(format string, values ...interface{}) error {
	return errors.New(fmt.Sprintf(format, values...))
}
