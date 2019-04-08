package str

import (
	"fmt"
	"testing"
)

func TestFloatDecimal(t *testing.T) {
	fmt.Println(FloatDecimalPrecise(0.01, 1))
	fmt.Println(FloatDecimalPrecise(1.1234567, 1))
	fmt.Println(FloatDecimalPrecise(1.1234567, 2))
	fmt.Println(FloatDecimalPrecise(1.1234567, 3))
	fmt.Println(FloatDecimalPrecise(1.1234567, 4))
	fmt.Println(FloatDecimalPrecise(1.1234567, 5))
	fmt.Println(FloatDecimalPrecise(1.1234567, 6))
}
