package uexcel

import (
	"fmt"
	"testing"
)

func TestX(t *testing.T) {
	for i := 0; i < 26; i++ {
		fmt.Println(X(i))
	}
}
