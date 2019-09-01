package uto

import (
	"encoding/json"
	"github.com/gosexy/to"
	"github.com/sirupsen/logrus"
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

func MapStringInterfaceToStringString(
	mInter map[string]interface{},
) (mStr map[string]string) {
	mStr = make(map[string]string)

	for key, vInter := range mInter {
		var strInter string
		bIntr, err := json.Marshal(vInter)
		if nil != err {
			logrus.Errorf("marshal interface value failed. error: %s.", err)
			strInter = to.String(vInter)
		} else {
			strInter = string(bIntr)
		}

		mStr[key] = strInter
	}

	return
}
