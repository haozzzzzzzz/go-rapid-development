package dyndb

import (
	"testing"

	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestIn(t *testing.T) {
	attributes := make(map[string]*dynamodb.AttributeValue)
	publisherIds := []interface{}{
		1, 2, 3,
	}

	condition, err := In("publisher_id", attributes, publisherIds)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(condition)
	fmt.Println(attributes)
}
