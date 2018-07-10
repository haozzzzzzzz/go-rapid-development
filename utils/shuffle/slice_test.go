package shuffle

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestShuffleSlice(t *testing.T) {
	rand.Seed(time.Now().Unix())

	strs := []interface{}{
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"10",
	}
	ShuffleSlice(strs)

	fmt.Println(strs)
}
