package str

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
)

func ToFloat(str string) (result float64, err error) {
	_, err = fmt.Sscanf(str, "%f", &result)
	if nil != err {
		logrus.Errorf("str to float failed. str: %s, err: %s", str, err)
		return
	}
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
