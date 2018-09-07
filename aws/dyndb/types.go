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

// 普通属性map可以转换成dynamodb的属性map
type AttributeValueKeyMap map[string]interface{}

func (m *AttributeValueKeyMap) Convert() (attributeValues map[string]*dynamodb.AttributeValue) {
	if *m == nil {
		return
	}

	attributeValues = make(map[string]*dynamodb.AttributeValue)

	for name, value := range *m {
		attributeValues[name] = KeyAttributeValue(value)
	}

	return
}
