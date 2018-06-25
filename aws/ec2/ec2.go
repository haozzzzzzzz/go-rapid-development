package ec2

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"
)

func GetEc2InstanceIdentityDocument(ses *session.Session) (doc ec2metadata.EC2InstanceIdentityDocument, err error) {
	client := ec2metadata.New(ses)
	doc, err = client.GetInstanceIdentityDocument()
	if nil != err {
		logrus.Errorf("Unable to read EC2 instance metadata: %v.", err)
		return
	}

	return
}
