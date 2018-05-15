package str

import (
	"fmt"
)

func ToFloat(str string) float64 {
	var result float64
	fmt.Sscanf(str, "%f", &result)
	return result
}
