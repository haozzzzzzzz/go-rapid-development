package dyndb

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"strings"
)

func ProjectionExpression(
	names []string,
) (expression string, expressionNames map[string]*string) {
	exprKeys := make([]string, 0)
	expressionNames = make(map[string]*string, 0)
	for _, name := range names {
		exprKey := fmt.Sprintf("#%s", name)
		exprKeys = append(exprKeys, exprKey)
		expressionNames[exprKey] = aws.String(name)
	}

	expression = strings.Join(exprKeys, ",")
	return
}
