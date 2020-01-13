package utime

import (
	"github.com/sirupsen/logrus"
	"time"
)

func ParseCommonDateFormatUnix(
	strTime string,
) (t time.Time, err error) {
	t, err = time.ParseInLocation(StrCommonDateFormat, strTime, time.UTC)
	if nil != err {
		logrus.Errorf("parse common-date-format time failed. error: %s.", err)
		return
	}
	return
}
