package str

import (
	"fmt"
	"github.com/xiam/to"
	"testing"
)

func TestToFloat(t *testing.T) {
	fmt.Println(ToFloat("1.1"))
}

func TestFloatDecimal(t *testing.T) {
	fmt.Println(FloatDecimalPrecise(0.01, 1))
	fmt.Println(FloatDecimalPrecise(1.1234567, 1))
	fmt.Println(FloatDecimalPrecise(1.1234567, 2))
	fmt.Println(FloatDecimalPrecise(1.1234567, 3))
	fmt.Println(FloatDecimalPrecise(1.1234567, 4))
	fmt.Println(FloatDecimalPrecise(1.1234567, 5))
	fmt.Println(FloatDecimalPrecise(1.1234567, 6))

	fmt.Println(to.Float64(FloatDecimalPrecise(1.299999, 2)))
}
