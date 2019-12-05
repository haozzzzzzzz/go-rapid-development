package utime

import (
	"github.com/sirupsen/logrus"
	"time"
)

// parse duration by hour:min:sec. 11:00:01, etc.
func ParseDurationByHMS(strHourMinSec string) (dur time.Duration, err error) {
	t, err := time.ParseInLocation("15:04:05", strHourMinSec, time.UTC)
	if nil != err {
		logrus.Errorf("parse formatted string time failed. error: %s.", err)
		return
	}

	dur = time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second
	return
}

func ParseDurSecondsByHMS(strHourMinSec string) (durSec int64, err error) {
	dur, err := ParseDurationByHMS(strHourMinSec)
	if nil != err {
		logrus.Errorf("parse duration by hms failed. error: %s.", err)
		return
	}

	durSec = int64(dur.Seconds())
	return
}

func ParseHMSByDurSeconds(durSec int64) (strHourMinSec string) {
	t := time.Unix(durSec, 0).In(time.UTC)
	strHourMinSec = t.Format("15:04:05")
	return
}
