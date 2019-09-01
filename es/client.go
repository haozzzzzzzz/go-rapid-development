package es

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
	"net/http"
)

func NewClient(
	address []string,
	transport http.RoundTripper,
) (client *elasticsearch.Client, err error) {
	client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: address,
		Transport: transport,
	})
	if nil != err {
		logrus.Errorf("new elasticsearch client failed. error: %s.", err)
		return
	}
	return
}
