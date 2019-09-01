package es_xray

import (
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/haozzzzzzzz/go-rapid-development/es"
	"github.com/sirupsen/logrus"
)

func NewClientWithXRay(
	address []string,
	newCheckerFunc func() es.RoundTripChecker,
) (client *elasticsearch.Client, err error) {
	transport := xray.RoundTripper(es.NewTransportCheckRoundTripper(es.ShortTimeoutTransport, newCheckerFunc)) // 3 RoundTripper stack up
	client, err = es.NewClient(address, transport)
	if nil != err {
		logrus.Errorf("new es client failed. error: %s.", err)
		return
	}

	return
}
