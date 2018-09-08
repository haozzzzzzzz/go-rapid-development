package dyndb

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

// 查询。dynamodb只会一次query返回总大小1m的数据，要获取所有数据，请使用QueryPage
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

func (m *Client) QueryPage(input *dynamodb.QueryInput, fn func(*dynamodb.QueryOutput, bool) bool) (err error) {
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
func (m *Client) BatchWriteItem(input *dynamodb.BatchWriteItemInput) (output *dynamodb.BatchWriteItemOutput, err error) {
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, input)
		defer func() {
			checker.After(err)
		}()
	}

	output, err = m.DB.BatchWriteItemWithContext(m.Ctx, input)
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
