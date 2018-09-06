package dyndb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gosexy/to"
)

func KeyAttributeValue(val interface{}) (attributeValue *dynamodb.AttributeValue) {
	attributeValue = &dynamodb.AttributeValue{}
	switch val.(type) {
	case string:
		attributeValue.S = aws.String(val.(string))
	default:
		attributeValue.N = aws.String(to.String(val))
	}

	return
}
