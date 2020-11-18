package dyndb

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

type CommandChecker interface {
	Before(client *Client, input interface{})
	After(err error)
}

type CommandCheckerMaker interface {
	NewChecker() CommandChecker
}

var (
	ErrNoRows = errors.New("dynamodb: no rows in result set")
)

type Client struct {
	DB                  *dynamodb.DynamoDB
	Ctx                 context.Context
	CommandCheckerMaker CommandCheckerMaker
}

func NewClient(db *dynamodb.DynamoDB, ctx context.Context) *Client {
	return &Client{
		DB:  db,
		Ctx: ctx,
	}
}

func (m *Client) CommandChecker() CommandChecker {
	if m.CommandCheckerMaker == nil {
		return nil
	}

	return m.CommandCheckerMaker.NewChecker()
}

func (m *Client) ListTables(input *dynamodb.ListTablesInput) (output *dynamodb.ListTablesOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.ListTablesWithContext(m.Ctx, input)
	return
}

func (m *Client) CreateTable(input *dynamodb.CreateTableInput) (output *dynamodb.CreateTableOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.CreateTable(input)
	return
}

func (m *Client) DescribeTable(input *dynamodb.DescribeTableInput) (output *dynamodb.DescribeTableOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.DescribeTable(input)
	return
}

func (m *Client) UpdateTable(input *dynamodb.UpdateTableInput) (output *dynamodb.UpdateTableOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.UpdateTableWithContext(m.Ctx, input)
	return
}

func (m *Client) TableExists(tableName string) (exists bool, err error) {
	_, err = m.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if nil != err {
		awsErr, ok := err.(awserr.Error)
		if ok == true && awsErr.Code() == dynamodb.ErrCodeResourceNotFoundException {
			err = nil
			exists = false

		} else {
			logrus.Errorf("describe table failed. error: %s.", err)
			return

		}

	} else {
		exists = true
	}

	return
}

// 根据offset检索limit个记录
// offset only allow string and number
func (m *Client) Scan(input *dynamodb.ScanInput) (output *dynamodb.ScanOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.ScanWithContext(m.Ctx, input)
	return
}

// 扫描
func (m *Client) ScanPage(input *dynamodb.ScanInput, fn func(output *dynamodb.ScanOutput, lastPage bool) (cont bool)) (err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	err = m.DB.ScanPagesWithContext(m.Ctx, input, fn)
	return
}

func (m *Client) Query(
	input *dynamodb.QueryInput,
) (
	output *dynamodb.QueryOutput,
	err error,
) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.QueryWithContext(m.Ctx, input)
	return
}

func (m *Client) QueryPage(input *dynamodb.QueryInput, fn func(output *dynamodb.QueryOutput, lastPage bool) (cont bool)) (err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	err = m.DB.QueryPagesWithContext(m.Ctx, input, fn)
	return
}

// 获取单个记录
func (m *Client) GetItem(input *dynamodb.GetItemInput) (output *dynamodb.GetItemOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.GetItemWithContext(m.Ctx, input)
	if nil != err {
		return
	}

	if output.Item == nil {
		err = ErrNoRows
		output = nil
	}

	return
}

func (m *Client) GetItemObj(input *dynamodb.GetItemInput, obj interface{}) (err error) {
	output, err := m.GetItem(input)
	if nil != err && err != ErrNoRows {
		logrus.Errorf("dyndb get item failed. error: %s.", err)
		return
	}

	if err == ErrNoRows {
		return
	}

	err = dynamodbattribute.UnmarshalMap(output.Item, obj)
	if nil != err {
		logrus.Errorf("dynamodbattribute unmarshal failed. error: %s.", err)
		return
	}
	return
}

// 批量获取
func (m *Client) BatchGetItem(input *dynamodb.BatchGetItemInput) (output *dynamodb.BatchGetItemOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.BatchGetItemWithContext(m.Ctx, input)
	return
}

// put item
func (m *Client) PutItemObj(tableName string, obj interface{}) (output *dynamodb.PutItemOutput, err error) {
	attributes, err := dynamodbattribute.MarshalMap(obj)
	if nil != err {
		logrus.Errorf("marshal attribute maps failed. error: %s.", err)
		return
	}

	output, err = m.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      attributes,
	})
	if nil != err {
		logrus.Errorf("put item failed. error: %s.", err)
		return
	}

	return
}

func (m *Client) PutItem(input *dynamodb.PutItemInput) (output *dynamodb.PutItemOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.PutItemWithContext(m.Ctx, input)
	return
}

// batch put item
// https://docs.aws.amazon.com/zh_cn/amazondynamodb/latest/APIReference/API_BatchWriteItem.html
func (m *Client) BatchWriteItem(input *dynamodb.BatchWriteItemInput) (err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	tryCount := 0
	for {
		output, errBatch := m.DB.BatchWriteItemWithContext(m.Ctx, input)
		err = errBatch
		if err != nil {
			logrus.Errorf("batch write item with context failed. error: %s", err)
			return
		}

		if output == nil {
			break
		}

		lenItem := len(output.UnprocessedItems)
		if lenItem <= 0 {
			break
		}

		// 重试
		// 指数回退算法
		tryCount++
		if tryCount > 24 {
			err = uerrors.Newf("batch write item request retry too many times. try_times: %s", tryCount)
			logrus.Error(err)
			break
		}

		waitTime := (math.Exp2(float64(tryCount/2)) + 1) * 10
		time.Sleep(time.Duration(waitTime) * time.Millisecond)
		logrus.Debugf("batch write item request retry. try_times: %d, request_count: %d, wait_time: %f ms", tryCount, lenItem, waitTime)

		input.RequestItems = output.UnprocessedItems // 需要重试
	}

	return
}

// 删除
func (m *Client) DeleteItem(input *dynamodb.DeleteItemInput) (output *dynamodb.DeleteItemOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.DeleteItemWithContext(m.Ctx, input)
	return
}

// 删除表单
func (m *Client) DeleteTable(input *dynamodb.DeleteTableInput) (output *dynamodb.DeleteTableOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.DeleteTableWithContext(m.Ctx, input)
	return
}

// 更新
func (m *Client) UpdateItem(
	input *dynamodb.UpdateItemInput,
) (output *dynamodb.UpdateItemOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.UpdateItemWithContext(m.Ctx, input)
	return
}

func (m *Client) TransactWriteItems(
	input *dynamodb.TransactWriteItemsInput,
) (output *dynamodb.TransactWriteItemsOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.TransactWriteItems(input)
	return
}
