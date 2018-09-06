package dyndb

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
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

func (m *Client) ListTables(offsetTableName string, limit int64) (output *dynamodb.ListTablesOutput, err error) {
	var inputExclusiveStartTableName *string
	if offsetTableName != "" {
		inputExclusiveStartTableName = aws.String(offsetTableName)
	}

	var inputLimit *int64
	if limit > 0 {
		inputLimit = aws.Int64(limit)
	}

	input := &dynamodb.ListTablesInput{
		ExclusiveStartTableName: inputExclusiveStartTableName,
		Limit: inputLimit,
	}

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

func (m *Client) DescribeTable(tableName string) (output *dynamodb.DescribeTableOutput, err error) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

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
func (m *Client) ScanByExclusive(tableName string, exclusiveMap map[string]interface{}, limit int64) (output *dynamodb.ScanOutput, err error) {
	exclusive := make(map[string]*dynamodb.AttributeValue)
	for name, value := range exclusiveMap {
		exclusive[name] = KeyAttributeValue(value)
	}

	if len(exclusive) == 0 {
		exclusive = nil
	}

	input := &dynamodb.ScanInput{
		TableName:         aws.String(tableName),
		Limit:             aws.Int64(limit),
		ExclusiveStartKey: exclusive,
	}

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
func (m *Client) QueryByExclusive(
	tableName string,
	partitionKey map[string]interface{}, // 分区键
	exclusive map[string]interface{}, // 主键
	limit int64,
) (
	output *dynamodb.QueryOutput,
	err error,
) {
	if len(partitionKey) != 1 {
		err = errors.New("need one partition key")
		return
	}

	var partitionKeyName string
	var partitionKeyValue interface{}
	for partitionKeyName, partitionKeyValue = range partitionKey {
		break
	}

	exclusiveStartKey := make(map[string]*dynamodb.AttributeValue)
	for name, value := range exclusive {
		exclusiveStartKey[name] = KeyAttributeValue(value)
	}

	if len(exclusiveStartKey) == 0 {
		exclusiveStartKey = nil
	}

	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("#partition_key=:partition_key"),
		ExpressionAttributeNames: map[string]*string{
			"#partition_key": aws.String(partitionKeyName),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":partition_key": KeyAttributeValue(partitionKeyValue),
		},
		ExclusiveStartKey: exclusiveStartKey,
		Limit:             aws.Int64(limit),
	}

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
func (m *Client) GetItem(tableName string, key map[string]interface{}) (output *dynamodb.GetItemOutput, err error) {
	attributeKey := make(map[string]*dynamodb.AttributeValue)
	for name, value := range key {
		attributeKey[name] = KeyAttributeValue(value)
	}

	if len(attributeKey) == 0 {
		err = ErrInvalidParams
		return
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       attributeKey,
	}

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

// 删除
func (m *Client) DeleteItem(tableName string, keyMap map[string]interface{}) (output *dynamodb.DeleteItemOutput, err error) {
	key := make(map[string]*dynamodb.AttributeValue)
	for name, value := range keyMap {
		key[name] = KeyAttributeValue(value)
	}

	if len(key) == 0 {
		err = ErrInvalidParams
		return
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}

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
	tableName string,
	keyMap map[string]interface{},
	updateExpression string,
	expressionAttributeNames map[string]*string,
	expressionAttributeValues map[string]*dynamodb.AttributeValue,
) (output *dynamodb.UpdateItemOutput, err error) {
	key := make(map[string]*dynamodb.AttributeValue)
	for name, value := range keyMap {
		key[name] = KeyAttributeValue(value)
	}

	if len(key) == 0 {
		err = ErrInvalidParams
		return
	}

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ExpressionAttributeNames:  expressionAttributeNames,
	}

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
