package region

import (
	"fmt"
	"testing"
)

func TestParseCommonDateFormat(t *testing.T) {
	tt, err := ParseCommonDateFormat("2020-03-24 00:00:00")
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(tt.Weekday())
}
