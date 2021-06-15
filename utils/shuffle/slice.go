package shuffle

import (
	"github.com/haozzzzzzzz/go-rapid-development/v2/utils/num"
)

// 打乱切片元素
func ShuffleSlice(rawSlice []interface{}) {
	lenSlice := len(rawSlice)
	if lenSlice == 0 {
		return
	}

	for i := 0; i < lenSlice-1; i++ {
		changeIndex := num.RandInt(i, lenSlice)
		if changeIndex == i {
			continue
		}

		temp := rawSlice[i]
		rawSlice[i] = rawSlice[changeIndex]
		rawSlice[changeIndex] = temp
	}

}

func ShuffleInt64Slice(rawSlice []int64) {
	lenSlice := len(rawSlice)
	if lenSlice == 0 {
		return
	}

	for i := 0; i < lenSlice-1; i++ {
		changeIndex := num.RandInt(i, lenSlice)
		if changeIndex == i {
			continue
		}

		temp := rawSlice[i]
		rawSlice[i] = rawSlice[changeIndex]
		rawSlice[changeIndex] = temp
	}

}

func ShuffleUint32Slice(rawSlice []uint32) {
	lenSlice := len(rawSlice)
	if lenSlice == 0 {
		return
	}

	for i := 0; i < lenSlice-1; i++ {
		changeIndex := num.RandInt(i, lenSlice)
		if changeIndex == i {
			continue
		}

		temp := rawSlice[i]
		rawSlice[i] = rawSlice[changeIndex]
		rawSlice[changeIndex] = temp
	}

}
