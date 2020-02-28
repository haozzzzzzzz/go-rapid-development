package str

import (
	"fmt"
	"testing"
)

func TestVersionOrdinal(t *testing.T) {
	res := VersionOrdinal("9.1.0") < VersionOrdinal("10")
	fmt.Println(res)
}
