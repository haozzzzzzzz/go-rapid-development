package uexcel

import (
	"fmt"
)

const NumChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func X(iX int) (x string) {
	if iX < 26 {
		x = string(NumChars[iX])
	}
	return
}

func Axis(x int, y int) (axis string) {
	return fmt.Sprintf("%s%d", X(x), y)
}
