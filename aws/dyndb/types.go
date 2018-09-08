package dyndb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

func KeyAttributeValue(val interface{}) (attributeValue *dynamodb.AttributeValue) {
	attributeValue, err := dynamodbattribute.Marshal(val)
	if nil != err {
		logrus.Errorf("dynamodbattribute marshal failed. error: %s.", err)
		attributeValue = nil
		return
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
