package num

import (
	"math/rand"
)

//[min, max)
func RandInt64(min, max int64) int64 {
	if min >= max {
		return max
	}

	return rand.Int63n(max-min) + min
}

//[min,max)
func RandInt(min, max int) int {
	if min >= max {
		return max
	}

	return rand.Intn(max-min) + min
}
