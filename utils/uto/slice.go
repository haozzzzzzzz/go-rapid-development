package uto

import (
	"github.com/gosexy/to"
)

func SliceStringToInt64(
	strSlice []string,
) (int64Slice []int64) {
	int64Slice = make([]int64, 0)
	for _, str := range strSlice {
		int64Slice = append(int64Slice, to.Int64(str))
	}
	return
}

func SliceInt64ToString(
	int64Slice []int64,
) (strSlice []string) {
	strSlice = make([]string, 0)
	for _, i := range int64Slice {
		strSlice = append(strSlice, to.String(i))
	}
	return
}
