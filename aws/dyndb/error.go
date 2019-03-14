package dyndb

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func IsErrConditionalCheckFailedException(
	err error,
) (isErr bool) {
	awsErr, ok := err.(awserr.Error)
	if ok && awsErr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
		isErr = true
	}
	return
}
