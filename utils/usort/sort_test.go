package usort

import (
	"fmt"
	"testing"
)

func TestInt64Slice(t *testing.T) {
	i := Int64Slice{
		8, 3, 9, 0, 2, 10, -1, 5, 4, 2, 1, 5, 6,
	}
	i.SortDesc()
	fmt.Println(i)
}
