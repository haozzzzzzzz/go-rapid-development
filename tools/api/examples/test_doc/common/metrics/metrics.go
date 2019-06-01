package metrics

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_doc/common/config"
)

var (
	MetricsNamespace string

	// ec2 instance id
	Ec2InstanceId string
)

func init() {
	var err error

	MetricsNamespace = config.ServiceConfig.MetricsNamespace
	Ec2InstanceId = config.AWSEc2InstanceIdentifyDocument.InstanceID
}
