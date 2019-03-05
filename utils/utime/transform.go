package utime

import "time"

// 一天的开始时间
func DayStartTime(t time.Time) (dayStartTime time.Time) {
	year, month, day := t.Date()
	dayStartTime = time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return
}

func DayStartTimeOffset(t time.Time, dayOffset int) (dayStartTime time.Time) {
	thisDayStartTime := DayStartTime(t)
	dayStartTime = thisDayStartTime.Add(24 * time.Duration(dayOffset) * time.Hour)
	return
}

// 一周的开始时间，从周一开始
func WeekStartTime(t time.Time) (weekStartTime time.Time) {
	dayStart := DayStartTime(t)

	dayStartUnix := dayStart.Unix()
	weekDay := dayStart.Weekday()

	// 本周周一
	var weekStartUnix int64
	var secondsPerDay int64 = 60 * 60 * 24
	offset := int64(int(weekDay)-1) * secondsPerDay
	weekStartUnix = dayStartUnix - offset
	weekStartTime = time.Unix(weekStartUnix, 0).In(t.Location())

	if weekDay == 0 { // 周日属于上一周
		weekStartTime = WeekStartTimeOffset(weekStartTime, -1)
	}

	return
}

func MonthStartTime(t time.Time) (startTime time.Time) {
	year, month, _ := t.Date()
	startTime = time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	return
}

func MonthStartTimeOffset(t time.Time, monthOffset int) (monthStartTime time.Time) {
	thisMonth := MonthStartTime(t)
	monthStartTime = thisMonth.AddDate(0, monthOffset, 0)
	return
}

// 前（后）几周的开始时间
func WeekStartTimeOffset(t time.Time, offset int) (thatWeekStartTime time.Time) {
	weekStartTime := WeekStartTime(t)
	var secondsWeekOffset = 60 * 60 * 24 * 7 * time.Duration(offset)
	thatWeekStartTime = weekStartTime.Add(secondsWeekOffset * time.Second)
	return
}

// 获取每周开始时间的标识
func WeekTimeIdentifyKey(t time.Time) string {
	return t.Format("w_2006_01_02")
}

func CommonDateFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func DateStringFormat(t time.Time) string {
	return t.Format("2006-01-02")
}
