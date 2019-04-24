package num

import (
	"fmt"
	"testing"
)

func TestIntersectIntRange(t *testing.T) {
	r1 := &IntRange{
		Min: 10,
		Max: 30,
	}

	r2 := &IntRange{
		Min: 20,
		Max: 21,
	}

	isIntersect := IntersectIntRange(r1, r2)
	fmt.Println(isIntersect)

}
