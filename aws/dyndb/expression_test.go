package dyndb

import (
	"fmt"
	"testing"
)

func TestProjectionExpression(t *testing.T) {
	expression, expressionNames := ProjectionExpression([]string{"exist", "day"})
	fmt.Println(expression)
	fmt.Println(expressionNames)
}
