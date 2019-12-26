package es_xray

import (
	"github.com/aws/aws-xray-sdk-go/xray"
	elasticsearch_v6 "github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/haozzzzzzzz/go-rapid-development/es"
	"github.com/sirupsen/logrus"
)

func NewClientWithXRay(
	address []string,
	newCheckerFunc func() es.IRoundTripChecker,
) (client *elasticsearch.Client, err error) {
	transport := xray.RoundTripper(es.NewTransportCheckRoundTripper(es.ShortTimeoutTransport, newCheckerFunc)) // 3 RoundTripper stack up
	client, err = es.NewClient(address, transport)
	if nil != err {
		logrus.Errorf("new es client failed. error: %s.", err)
		return
	}

	return
}

func NewClientV6WithXRay(
	address []string,
	newCheckerFunc func() es.IRoundTripChecker,
) (client *elasticsearch_v6.Client, err error) {
	transport := xray.RoundTripper(es.NewTransportCheckRoundTripper(es.ShortTimeoutTransport, newCheckerFunc)) // 3 RoundTripper stack up
	client, err = es.NewClientV6(address, transport)
	if nil != err {
		logrus.Errorf("new es client failed. error: %s.", err)
		return
	}
	return
}

func NewClientV6NoTimeoutWithXRay(
	address []string,
	newCheckerFunc func() es.IRoundTripChecker,
) (client *elasticsearch_v6.Client, err error) {
	transport := xray.RoundTripper(es.NewTransportCheckRoundTripper(es.NoTimeoutTransport, newCheckerFunc)) // 3 RoundTripper stack up
	client, err = es.NewClientV6(address, transport)
	if nil != err {
		logrus.Errorf("new es client failed. error: %s.", err)
		return
	}
	return
}
