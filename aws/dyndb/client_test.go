package dyndb

import (
	"os"
	"testing"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/haozzzzzzzz/go-rapid-development/aws/xray"
	"github.com/sirupsen/logrus"
)

var db *dynamodb.DynamoDB

func TestMain(m *testing.M) {
	var err error
	db, err = NewDynamoDB(&ClientConfigFormat{
		Endpoint: "http://localhost:8000",
		Region:   "ap-south-1",
		Credentials: &ClientCredentials{
			ID:     "a",
			Secret: "b",
			Token:  "c",
		},
	})
	if nil != err {
		logrus.Errorf("new dynamodb failed. error: %s.", err)
		return
	}
	os.Exit(m.Run())
}

func TestClient_TableExists(t *testing.T) {
	var err error
	ctx, _, cancel := xray.NewBackgroundContext("test")
	defer func() {
		cancel(err)
	}()

	client := Client{
		DB:  db,
		Ctx: ctx,
	}

	exists, err := client.TableExists("movies")
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(exists)
}

func TestClient_GetItem(t *testing.T) {
	var err error
	ctx, _, cancel := xray.NewBackgroundContext("test")
	defer func() {
		cancel(err)
	}()

	client := NewClient(db, ctx)
	_, err = client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("dev_simple_store"),
		Key: map[string]*dynamodb.AttributeValue{
			"hk_type": {
				S: aws.String("commodity_order_excel"),
			},
			"rk_id": {
				S: aws.String("2"),
			},
		},
	})
	if nil != err {
		t.Error(err)
		return
	}
}
