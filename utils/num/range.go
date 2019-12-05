package num

import (
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
)

// 全闭合的
type IntRange struct {
	Min int64
	Max int64
}

func NewIntRange(min int64, max int64) (intRange *IntRange, err error) {
	if max < min {
		err = uerrors.Newf("int range max should be gte min")
		return
	}
	intRange = &IntRange{
		Min: min,
		Max: max,
	}
	return
}

func (m *IntRange) IsBetween(val int64) (result bool) {
	if m.Min <= val && val <= m.Max {
		result = true
	}

	return
}

func IntersectIntRange(
	r1 *IntRange,
	r2 *IntRange,
) (isIntersect bool) {
	if r1.IsBetween(r2.Min) ||
		r1.IsBetween(r2.Max) ||
		r2.IsBetween(r1.Min) ||
		r2.IsBetween(r1.Max) {
		isIntersect = true
		return
	}
	return
}

func IntRangeIntersectRange(
	r1 *IntRange,
	r2 *IntRange,
) (isIntersect bool, sub *IntRange) {
	isIntersect = IntersectIntRange(r1, r2)
	if !isIntersect {
		return
	}

	sub = &IntRange{}
	sub.Min = r1.Min
	if sub.Min < r2.Min { // 取min大值
		sub.Min = r2.Min
	}

	sub.Max = r1.Max
	if sub.Max > r2.Max { // 取max小值
		sub.Max = r2.Max
	}

	return
}
