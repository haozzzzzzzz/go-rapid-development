package ubuffer

import (
	"bytes"

	"github.com/sirupsen/logrus"
)

type OffsetByteBuffer struct {
	bytes.Buffer
}

func (m *OffsetByteBuffer) WriteAt(p []byte, off int64) (n int, err error) {
	n, err = m.Write(p)
	if nil != err {
		logrus.Warnf("write bytes failed. %s", err)
		return
	}
	return
}

func NewOffsetByteBuffer(buf []byte) *OffsetByteBuffer {
	return &OffsetByteBuffer{
		Buffer: *bytes.NewBuffer(buf),
	}
}
