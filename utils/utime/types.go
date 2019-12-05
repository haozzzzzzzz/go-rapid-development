package utime

import (
	"errors"
	"github.com/haozzzzzzzz/go-rapid-development/utils/num"
	"github.com/sirupsen/logrus"
)

// [start_time, end_time)
type TimeRange struct {
	StartTime int64 `json:"start_time" yaml:"start_time" validate:"required"`
	EndTime   int64 `json:"end_time" yaml:"end_time" validate:"required"`
}

func NewTimeRange(startTime, endTime int64) (tr *TimeRange, err error) {
	if startTime >= endTime {
		err = errors.New("time range start_time should less than end_time")
		return
	}
	tr = &TimeRange{
		StartTime: startTime,
		EndTime:   endTime,
	}
	return
}

/*
0: t < start_time
1: start_time <= t < end_time
2: t >= end_time
*/
const TimeRangeStateNotStart = 0
const TimeRangeStateStarted = 1
const TimeRangeStateOver = 2

func (m *TimeRange) TimeState(t int64) (state int) {
	state = TimeRangeStateNotStart
	if t >= m.StartTime {
		if t < m.EndTime {
			state = TimeRangeStateStarted
		} else {
			state = TimeRangeStateOver
		}
	}
	return
}

//func (m *TimeRange) ToIntRange() *num.IntRange {
//	return &num.IntRange{
//		Min: m.StartTime,
//		Max: m.EndTime,
//	}
//}

//func NewTimeRangeFromIntRange(r *num.IntRange) *TimeRange {
//	return &TimeRange{
//		StartTime: r.Min,
//		EndTime:   r.Max,
//	}
//}

// Deprecated: 结束时间与其他的开始时间重合也是重合
func IntersectTimeRange(
	r1 *TimeRange,
	r2 *TimeRange,
) (isIntersect bool, sub *TimeRange) {
	r1IntRange := &num.IntRange{
		Min: r1.StartTime,
		Max: r1.EndTime,
	}

	r2IntRange := &num.IntRange{
		Min: r2.StartTime,
		Max: r2.EndTime,
	}

	isIntersect, ir := num.IntRangeIntersectRange(
		r1IntRange,
		r2IntRange,
	)

	if !isIntersect {
		return
	}

	sub = &TimeRange{
		StartTime: ir.Min,
		EndTime:   ir.Max,
	}

	return
}

// 结束时间与其他的开始时间重合不算重合
func IntersectTimeRangeStrict(
	r1 *TimeRange,
	r2 *TimeRange,
) (isIntersect bool, sub *TimeRange, err error) {

	// 左闭右开转换为全闭合
	r1IntRange, err := num.NewIntRange(r1.StartTime, r1.EndTime-1)
	if nil != err {
		logrus.Errorf("new int range failed. error: %s.", err)
		return
	}

	r2IntRange, err := num.NewIntRange(r2.StartTime, r2.EndTime-1)
	if nil != err {
		logrus.Errorf("new int range failed. error: %s.", err)
		return
	}

	isIntersect, ir := num.IntRangeIntersectRange(
		r1IntRange,
		r2IntRange,
	)

	if !isIntersect {
		return
	}

	sub = &TimeRange{
		StartTime: ir.Min,
		EndTime:   ir.Max + 1,
	}

	return
}
