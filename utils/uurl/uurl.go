package uurl

import (
	"fmt"
)

type TestUint32 uint32

func ToQueryPairs(paramsMap map[string]string) (strQueryPairs []string) {
	strQueryPairs = make([]string, 0)
	for key, value := range paramsMap {
		strQueryPairs = append(strQueryPairs, fmt.Sprintf("%s=%s", key, value))
	}
	return
}
