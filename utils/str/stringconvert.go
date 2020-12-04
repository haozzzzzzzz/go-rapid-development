package str

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xiam/to"
	"math"
)

func ToFloat(str string) (result float64, err error) {
	_, err = fmt.Sscanf(str, "%f", &result)
	if nil != err {
		logrus.Errorf("str to float failed. str: %s, err: %s", str, err)
		return
	}
	return
}

// 精确截断小数位
func FloatDecimalPrecise(n float64, bit uint32) (strDecimal string) {
	bitNum := math.Pow10(int(bit))
	n = math.Trunc(n*float64(bitNum)) / bitNum
	strDecimal = FloatDecimalRound(n, bit)
	return
}

func FloatDecimal(n float64, bit uint32) (newVal float64) {
	return to.Float64(FloatDecimalPrecise(n, bit))
}

// 四舍五入
func FloatDecimalRound(n float64, bit uint32) (strDecimal string) {
	strBitFormat := fmt.Sprintf("%%.%df", bit)
	strDecimal = fmt.Sprintf(strBitFormat, n)
	return
}
