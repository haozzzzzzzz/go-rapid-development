package gofmt

import (
	"errors"
	"fmt"
	"github.com/haozzzzzzzz/go-rapid-development/cmd"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/sirupsen/logrus"
	"go/format"
	"os"
	"os/exec"
)

func CmdGoFmt(path string) (err error) {
	goRoot := os.Getenv("GOROOT")
	goBin := ""
	if goRoot != "" {
		goBin = fmt.Sprintf("%s/bin/go", goRoot)
	} else {
		goBin, err = exec.LookPath("go")
		if nil != err {
			logrus.Errorf("look for go bin failed. error: %s.", err)
			return
		}
	}

	if goBin == "" {
		err = errors.New("go bin not specified")
		return
	}

	code, err := cmd.RunCommand("", goBin, "fmt", path)
	if nil != err {
		logrus.Errorf("run go fmt cmd failed. error: %s.", err)
		return
	}

	if code != 0 {
		err = uerrors.Newf("cmd exit abnormally. code: %d", code)
		return
	}

	return
}

func StrGoFmt(buf string) (strFormatted string, err error) {
	defer func() {
		if iRec := recover(); iRec != nil {
			err = uerrors.Newf("go fmt string panic: %s", err)
		}
	}()
	formatted, err := format.Source([]byte(buf))
	if err != nil {
		panic(fmt.Errorf("%s\nOriginal code:\n%s", err.Error(), buf))
	}
	strFormatted = string(formatted)
	return
}
