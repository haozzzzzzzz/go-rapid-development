package utime

import (
	"github.com/haozzzzzzzz/go-rapid-development/utils/num"
)

type TimeRange struct {
	StartTime int64 `json:"start_time" yaml:"start_time" validate:"required"`
	EndTime   int64 `json:"end_time" yaml:"end_time" validate:"required"`
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

func (m *TimeRange) ToIntRange() *num.IntRange {
	return &num.IntRange{
		Min: m.StartTime,
		Max: m.EndTime,
	}
}

func NewTimeRangeFromIntRange(r *num.IntRange) *TimeRange {
	return &TimeRange{
		StartTime: r.Min,
		EndTime:   r.Max,
	}
}

func IntersectTimeRange(
	r1 *TimeRange,
	r2 *TimeRange,
) (isIntersect bool, sub *TimeRange) {
	isIntersect, ir := num.IntRangeIntersectRange(
		r1.ToIntRange(),
		r2.ToIntRange(),
	)

	if !isIntersect {
		return
	}

	sub = NewTimeRangeFromIntRange(ir)

	return
}
