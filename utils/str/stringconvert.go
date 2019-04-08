package str

import (
	"fmt"
	"math"
)

func ToFloat(str string) float64 {
	var result float64
	fmt.Sscanf(str, "%f", &result)
	return result
}

func FloatDecimalPrecise(n float64, bit uint32) (strDecimal string) {
	bitNum := math.Pow10(int(bit))
	n = math.Trunc(n*float64(bitNum)) / bitNum
	strDecimal = FloatDecimalRound(n, bit)
	return
}

func FloatDecimalRound(n float64, bit uint32) (strDecimal string) {
	strBitFormat := fmt.Sprintf("%%.%df", bit)
	strDecimal = fmt.Sprintf(strBitFormat, n)
	return
}
