package time

import "time"

// 一天的开始时间
func DayStartTime(t time.Time) (dayStartTime time.Time, err error) {
	strDayStartTime := t.Format("2006-01-02 00:00:00 -0700 MST")
	dayStartTime, err = time.Parse("2006-01-02 15:04:05 -0700 MST", strDayStartTime)
	return
}

func TodayStartTime() (today time.Time, err error) {
	today, err = DayStartTime(time.Now())
	return
}

// 一周的开始时间，从周一开始
func WeekStartTime(t time.Time) (weekStartTime time.Time, err error) {
	dayStart, err := DayStartTime(t)
	if err != nil {
		return
	}

	dayStartUnix := dayStart.Unix()
	weekDay := dayStart.Weekday()

	// 本周周一
	var weekStartUnix int64
	var secondsPerDay int64 = 60 * 60 * 24
	offset := int64(int(weekDay)-1) * secondsPerDay
	weekStartUnix = dayStartUnix - offset
	weekStartTime = time.Unix(weekStartUnix, 0)

	if weekDay == 0 { // 周日属于上一周
		weekStartTime = OffsetWeekStartTime(weekStartTime, -1)
	}

	return
}

// 前（后）几周的开始时间
func OffsetWeekStartTime(weekStartTime time.Time, offset int) (thatWeekStartTime time.Time) {
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
