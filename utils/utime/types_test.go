package utime

import (
	"fmt"
	"testing"
)

func TestIntersectTimeRange(t *testing.T) {
	isIntersect, intersect := IntersectTimeRange(&TimeRange{
		StartTime: 1,
		EndTime:   2,
	}, &TimeRange{
		StartTime: 2,
		EndTime:   3,
	})
	fmt.Println(isIntersect, intersect)
}

func TestIntersectTimeRangeStrict(t *testing.T) {
	isIntersect, intersect, err := IntersectTimeRangeStrict(&TimeRange{
		StartTime: 1,
		EndTime:   2,
	}, &TimeRange{
		StartTime: 2,
		EndTime:   3,
	})
	if nil != err {
		t.Error(err)
		return
	}
	fmt.Println(isIntersect, intersect)
}
