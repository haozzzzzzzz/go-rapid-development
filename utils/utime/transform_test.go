package utime

import (
	"fmt"
	"testing"
	"time"
)

func TestDayStartTime(t *testing.T) {
	dayStartTime := DayStartTime(time.Now())
	fmt.Println(dayStartTime)
}
