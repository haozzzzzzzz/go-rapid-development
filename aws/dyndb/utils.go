package dyndb

import (
	"fmt"

	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

func In(
	prefix string, // expression name prefix
	outAttributeValueMap map[string]*dynamodb.AttributeValue,
	items []interface{},
) (condition string, err error) {
	names := make([]string, 0)
	for index, item := range items {
		name := fmt.Sprintf(":%s_in_%d", prefix, index)
		names = append(names, name)

		outAttributeValueMap[name], err = dynamodbattribute.Marshal(item)
		if nil != err {
			logrus.Errorf("marshal attribute value failed. arg: %v. error: %s.", item, err)
			return
		}
	}

	condition = strings.Join(names, ",")
	return
}
