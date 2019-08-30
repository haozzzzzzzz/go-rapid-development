package es

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
)

func NewClient(
	address []string,
	newCheckerFunc func() RoundTripChecker,
) (client *elasticsearch.Client, err error) {
	client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: address,
		Transport: NewTransportCheck(
			ShortTimeoutTransport,
			newCheckerFunc,
		),
	})
	if nil != err {
		logrus.Errorf("new elasticsearch client failed. error: %s.", err)
		return
	}

	return
}
