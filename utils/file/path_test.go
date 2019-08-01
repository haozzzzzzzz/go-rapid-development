package file

import (
	"fmt"
	"testing"
)

func TestParentDir(t *testing.T) {
	fmt.Println(ParentDir("/"))
}
