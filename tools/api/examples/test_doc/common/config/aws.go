package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/validator"
	"github.com/haozzzzzzzz/go-rapid-development/aws/ec2"
	"github.com/sirupsen/logrus"
)

type AWSConfigFormat struct {
	Region string `yaml:"region" validate:"required"`
}

var AWSConfig *AWSConfigFormat
var AWSSession *session.Session
var AWSEc2InstanceIdentifyDocument ec2metadata.EC2InstanceIdentityDocument

func CheckAWSConfig() {
	if AWSConfig != nil {
		return
	}

	CheckConsulConfig()
	AWSConfig = &AWSConfigFormat{}
	var err error
	err = GetConsulClient().GetYaml(ConsulConfig.KeyPrefix+"/aws.yaml", AWSConfig)
	if nil != err {
		logrus.Fatalf("read aws config file failed. %s", err)
		return
	}

	err = validator.New().Struct(AWSConfig)
	if nil != err {
		logrus.Fatalf("validate aws config file failed. %s", err)
		return
	}

	awsConfig := &aws.Config{}
	awsConfig.Region = aws.String(AWSConfig.Region)
	AWSSession, err = session.NewSession(awsConfig)
	if nil != err {
		logrus.Fatalf("new aws session failed. %s", err)
		return
	}

	AWSEc2InstanceIdentifyDocument, err = ec2.GetEc2InstanceIdentityDocument(AWSSession)
	if nil != err {
		logrus.Errorf("get ec2 instance identify document failed. %s.", err)
		return
	}

}
