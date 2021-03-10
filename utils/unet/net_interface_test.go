package unet

import (
	"fmt"
	"testing"
)

func TestGetInterfaceIP(t *testing.T) {
	ip, err := GetInterfaceIP("en0")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(ip)
}
