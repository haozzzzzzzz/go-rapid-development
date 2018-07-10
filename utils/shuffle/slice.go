package shuffle

import (
	"github.com/haozzzzzzzz/go-rapid-development/utils/num"
)

// 打乱切片元素
func ShuffleSlice(rawSlice []interface{}) {
	lenSlice := len(rawSlice)
	if lenSlice == 0 {
		return
	}

	for i := 0; i < lenSlice-1; i++ {
		changeIndex := num.RandInt(i+1, lenSlice)
		temp := rawSlice[i]
		rawSlice[i] = rawSlice[changeIndex]
		rawSlice[changeIndex] = temp
	}

}
