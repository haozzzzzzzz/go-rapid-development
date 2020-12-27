package region

import (
	"fmt"
	"testing"
	"time"
)

func TestParseCommonDateFormat(t *testing.T) {
	tt, err := ParseCommonDateFormat("2020-03-24 00:00:00")
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(tt.Weekday())
}

func TestAppRegionHourStartTime(t *testing.T) {
	SetAppRegionTimezone(IndianStandardTimezone)
	now := time.Now()
	hourTime := AppRegionHourStartTime(now)
	fmt.Println(hourTime)
	h24 := AppRegionHourStartTimeOffset(now, -25)
	fmt.Println(h24)
}
