package types

import (
	"github.com/aws/aws-sdk-go/aws"
)

func StringNilIfEmpty(str string) *string {
	if str == "" {
		return nil
	}

	return aws.String(str)
}
