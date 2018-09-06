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

var clientDynamoDB *dynamodb.DynamoDB

func TestMain(m *testing.M) {
	var err error
	clientDynamoDB, err = NewDynamoDB(&ClientConfigFormat{
		Region:   "ap-south-1",
		Endpoint: "http://localhost:8000",
		Credentials: &ClientCredentials{
			ID:     "AKID",
			Secret: "SECRET",
			Token:  "TOKEN",
		},
	})
	if nil != err {
		logrus.Fatal(err)
		return
	}

	os.Exit(m.Run())
}

func TestClient_ListTables(t *testing.T) {
	var err error
	ctx, _, cancel := xray.NewBackgroundContext("test")
	defer func() {
		cancel(err)
	}()

	client := &Client{
		DB:  clientDynamoDB,
		Ctx: ctx,
	}

	output, err := client.ListTables("", 1)
	if nil != err {
		t.Error(err)
		return
	}
	fmt.Println(output)

}

func TestClient_DescribeTable(t *testing.T) {
	var err error
	ctx, _, cancel := xray.NewBackgroundContext("test")
	defer func() {
		cancel(err)
	}()

	client := &Client{
		DB:  clientDynamoDB,
		Ctx: ctx,
	}

	output, err := client.DescribeTable("movies")
	if nil != err {
		logrus.Errorf("describe table failed. error: %s.", err)
		return
	}

	fmt.Println(output)
}

func TestClient_ScanByMapOffset(t *testing.T) {
	var err error
	ctx, _, cancel := xray.NewBackgroundContext("test")
	defer func() {
		cancel(err)
	}()

	client := &Client{
		DB:  clientDynamoDB,
		Ctx: ctx,
	}

	output, err := client.ScanByExclusive("movies", map[string]interface{}{
		"year":  2015,
		"title": "The Big New Movie 4",
	}, 1)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(output)
}

func TestClient_QueryByExclusive(t *testing.T) {
	var err error
	ctx, _, cancel := xray.NewBackgroundContext("test")
	defer func() {
		cancel(err)
	}()

	client := &Client{
		DB:  clientDynamoDB,
		Ctx: ctx,
	}

	output, err := client.QueryByExclusive("movies", map[string]interface{}{
		"year": 2015,
	},
		map[string]interface{}{
			"year":  2015,
			"title": "The Big New Movie 2",
		},
		1)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(output)
}

func TestClient_GetItem(t *testing.T) {
	var err error
	ctx, _, cancel := xray.NewBackgroundContext("test")
	defer func() {
		cancel(err)
	}()

	client := &Client{
		DB:  clientDynamoDB,
		Ctx: ctx,
	}

	output, err := client.GetItem("movies", map[string]interface{}{
		"year":  2015,
		"title": "The Big New Movie",
	})
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(output)
}

func TestClient_DeleteItem(t *testing.T) {
	var err error
	ctx, _, cancel := xray.NewBackgroundContext("test")
	defer func() {
		cancel(err)
	}()

	client := &Client{
		DB:  clientDynamoDB,
		Ctx: ctx,
	}

	output, err := client.DeleteItem("movies", map[string]interface{}{
		"year":  2015,
		"title": "The Big New Movie 2",
	})
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(output)
}

func TestClient_UpdateItem(t *testing.T) {
	var err error
	ctx, _, cancel := xray.NewBackgroundContext("test")
	defer func() {
		cancel(err)
	}()

	client := &Client{
		DB:  clientDynamoDB,
		Ctx: ctx,
	}

	output, err := client.UpdateItem("movies",
		map[string]interface{}{
			"year":  2015,
			"title": "The Big New Movie 5",
		},
		"SET info.rating = :r", nil, map[string]*dynamodb.AttributeValue{
			":r": {
				N: aws.String("6"),
			},
		})
	if nil != err {
		t.Error(err)
		return
	}
	fmt.Println(output)
}
