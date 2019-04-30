package num

type IntRange struct {
	Min int64
	Max int64
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