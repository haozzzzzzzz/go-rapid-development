package uio

import (
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
)

func ReadAllAndClose(reader io.ReadCloser) (bContent []byte, err error) {
	defer func() {
		errClose := reader.Close()
		if errClose != nil {
			logrus.Errorf("close reader failed. error: %s.", err)

			if err == nil {
				err = errClose
			}
		}
	}()

	bContent, err = ioutil.ReadAll(reader)
	if nil != err {
		logrus.Errorf("read all from ReadCloser failed. error: %s.", err)
		return
	}

	return
}
