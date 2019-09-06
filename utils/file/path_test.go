package file

import (
	"fmt"
	"testing"
)

func TestParentDir(t *testing.T) {
	fmt.Println(ParentDir("/"))
}

func TestDirHasExtFile(t *testing.T) {
	hasExt, err := DirHasExtFile("/Users/hao/Documents/Projects/Github/go-rapid-development/utils/file", ".go")
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(hasExt)
}
