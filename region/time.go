package region

import (
	"github.com/haozzzzzzzz/go-rapid-development/utils/utime"
	"github.com/sirupsen/logrus"
	"time"
)

// https://baike.baidu.com/item/印度标准时间/1313539?fr=aladdin
var IndianStandardTimezone = time.FixedZone("IST", 5*3600+1800) // +5:30

// https://24timezones.com/time-zone/wib
var WesternIndonesiaTimezone = time.FixedZone("WIB", 7*3600) // +7:00

var AppRegionTimezone = IndianStandardTimezone // 默认是印度时间

// set timezone
func SetAppRegionTimezone(loc *time.Location) {
	AppRegionTimezone = loc
}

func AppRegionNow() (t time.Time) {
	return time.Now().In(AppRegionTimezone)
}

func AppRegionNowUnix() int64 {
	return AppRegionNow().Unix()
}

func AppRegionUnix(sec int64, nsec int64) (t time.Time) {
	return time.Unix(sec, nsec).In(AppRegionTimezone)
}

func AppRegionDayStartTime(t time.Time) (dayStartTime time.Time) {
	year, month, day := t.In(AppRegionTimezone).Date()
	dayStartTime = time.Date(year, month, day, 0, 0, 0, 0, AppRegionTimezone)
	return
}

func AppRegionTodayStartTime() (today time.Time) {
	return AppRegionDayStartTime(AppRegionNow())
}

func AppRegionTomorrowStartTime() (tomorrow time.Time) {
	return AppRegionDayStartTimeOffset(AppRegionNow(), 1)
}

func AppRegionYesterdayStartTime() (yesterday time.Time) {
	return AppRegionDayStartTimeOffset(AppRegionNow(), -1)
}

func AppRegionCurrentDayRangeTime(t time.Time) (startTime time.Time, endTime time.Time) {
	startTime = AppRegionDayStartTime(t)
	endTime = AppRegionDayStartTimeOffset(startTime, 1)
	return
}

func AppRegionDayStartTimeOffset(t time.Time, dayOffset int) (dayStartTime time.Time) {
	thisDayStartTime := AppRegionDayStartTime(t)
	dayStartTime = thisDayStartTime.Add(24 * time.Hour * time.Duration(dayOffset))
	return
}

func AppRegionMonthStartTime(t time.Time) (startTime time.Time) {
	year, month, _ := t.In(AppRegionTimezone).Date()
	startTime = time.Date(year, month, 1, 0, 0, 0, 0, AppRegionTimezone)
	return
}

func AppRegionWeekStartTime(t time.Time) (startTime time.Time) {
	startTime = utime.WeekStartTime(t.In(AppRegionTimezone))
	return
}

func AppRegionWeekStartTimeOffset(t time.Time, offset int) (startTime time.Time) {
	startTime = utime.WeekStartTimeOffset(t.In(AppRegionTimezone), offset)
	return
}

func AppRegionMonthStartTimeOffset(t time.Time, monthOffset int) (monthStartTime time.Time) {
	thisMonth := AppRegionMonthStartTime(t)
	monthStartTime = thisMonth.AddDate(0, monthOffset, 0)
	return
}

func AppRegionParseDayTime(strDayTime string) (dayStartTime time.Time, err error) {
	dayStartTime, err = time.ParseInLocation("2006-01-02", strDayTime, AppRegionTimezone)
	if nil != err {
		logrus.Errorf("parse layout failed. error: %s.", err)
		return
	}
	return
}

// 当前时间距离每日特定时刻的时间间隔
// 用于计算每日时刻间隔
func AppRegionIndianMomentDuration(
	curTime time.Time, // 当前传入的时间
	dailyMomentDur time.Duration, // 每日的距离开始时间的间隔
) (distanceDur time.Duration) { // 当前时间距离下次时刻的时间间隔
	startTime := AppRegionDayStartTime(curTime)
	moment := startTime.Add(dailyMomentDur)

	if moment.After(curTime) { // 当天还没到
		distanceDur = moment.Sub(curTime)
		return
	}

	// 明天
	nextStartTime := AppRegionDayStartTimeOffset(curTime, 1)
	nextMoment := nextStartTime.Add(dailyMomentDur)

	distanceDur = nextMoment.Sub(curTime)

	return
}

func ParseCommonDateFormat(strTime string) (time.Time, error) {
	return time.ParseInLocation(utime.StrCommonDateFormat, strTime, AppRegionTimezone)
}

func AppRegionHourStartTime(t time.Time) (hourStartTime time.Time) {
	t = t.In(AppRegionTimezone)
	hourStartTime = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, AppRegionTimezone)
	return
}

func AppRegionHourStartTimeUnix(t time.Time) (hourUnix int64) {
	return AppRegionHourStartTime(t).Unix()
}

func AppRegionHourStartTimeOffset(t time.Time, hourOffset int) (hourStartTime time.Time) {
	curHourTime := AppRegionHourStartTime(t)
	hourStartTime = curHourTime.Add(time.Duration(hourOffset) * time.Hour)
	return
}
