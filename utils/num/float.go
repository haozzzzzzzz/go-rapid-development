package num

import (
	"fmt"
	"github.com/gosexy/to"
)

func FloatDecimal(origin float64, decimal int) float64 {
	preFmt := "%." + to.String(decimal) + "f"
	return to.Float64(fmt.Sprintf(preFmt, origin))
}
