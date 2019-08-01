package mod

import (
	"fmt"
	"testing"
)

func TestFindGoMod(t *testing.T) {
	curDir := "/Users/hao/Documents/Projects/XunLei/video_buddy_user_device/rpc/api/user_device"
	found, modName, modDir := FindGoMod(curDir)
	fmt.Println(found)
	fmt.Println(modName)
	fmt.Println(modDir)
}
