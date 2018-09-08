package dyndb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/haozzzzzzzz/go-rapid-development/aws/types"
	"github.com/sirupsen/logrus"
)

type ClientCredentials struct {
	ID     string `json:"id" yaml:"id"`
	Secret string `json:"secret" yaml:"secret"`
	Token  string `json:"token" yaml:"token"`
}

type ClientConfigFormat struct {
	Region      string             `json:"region" yaml:"region"`
	Credentials *ClientCredentials `json:"credentials" yaml:"credentials"`
	Endpoint    string             `json:"endpoint" yaml:"endpoint"`
}

func NewDynamoDB(config *ClientConfigFormat) (client *dynamodb.DynamoDB, err error) {
	var sessionCredentials *credentials.Credentials
	if config.Credentials != nil {
		sessionCredentials = credentials.NewStaticCredentials(config.Credentials.ID, config.Credentials.Secret, config.Credentials.Token)
	}

	ses, err := session.NewSession(&aws.Config{
		Region:      types.StringNilIfEmpty(config.Region),
		Endpoint:    types.StringNilIfEmpty(config.Endpoint),
		Credentials: sessionCredentials,
	})
	if nil != err {
		logrus.Errorf("new session failed. error: %s.", err)
		return
	}

	client = dynamodb.New(ses)

	return
}
