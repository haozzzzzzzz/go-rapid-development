package dyndb

import (
	"context"

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

var ErrNoItems = errors.New("dynamodb: no items in result set")
var ErrInvalidParams = errors.New("dynamodb: invalid params")

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
	if nil != err {
		logrus.Errorf("dynamodb list tables failed. error: %s.", err)
		return
	}

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
	if nil != err {
		logrus.Errorf("describe table failed. error: %s.", err)
		return
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
	if nil != err {
		logrus.Errorf("dynamodb scan failed. error: %s.", err)
		return
	}

	return
}

// 查询
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
	if nil != err {
		logrus.Errorf("dynamodb query failed. error: %s.", err)
		return
	}
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
		logrus.Errorf("dynamodb get item failed. error: %s.", err)
		return
	}

	if output.Item == nil {
		err = ErrNoItems
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
	if nil != err {
		logrus.Errorf("dynamodb batch get item failed. error: %s.", err)
		return
	}

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
	if nil != err {
		logrus.Errorf("dynamodb put item failed. error: %s.", err)
		return
	}

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
	if nil != err {
		logrus.Errorf("dynamodb batch write item failed. error: %s.", err)
		return
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
	if nil != err {
		logrus.Errorf("dynamodb delete item failed. error: %s.", err)
		return
	}

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
	if nil != err {
		logrus.Errorf("dynamodb update item failed. error: %s.", err)
		return
	}

	return
}
