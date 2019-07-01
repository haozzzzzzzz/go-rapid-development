package utime

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
