package utime

import (
	"fmt"
	"testing"
)

func TestParseDurSecondsByHMS(t *testing.T) {
	durSec, err := ParseDurSecondsByHMS("10:51:59")
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(durSec)
}

func TestParseHMSByDurSeconds(t *testing.T) {
	hms := ParseHMSByDurSeconds(3900)
	fmt.Println(hms)
}
