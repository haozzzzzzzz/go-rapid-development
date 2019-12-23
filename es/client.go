package es

import (
	elasticsearch_v6 "github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
	"net/http"
)

// client v7
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

// client v6
func NewClientV6(
	address []string,
	transport http.RoundTripper,
) (client *elasticsearch_v6.Client, err error) {
	client, err = elasticsearch_v6.NewClient(elasticsearch_v6.Config{
		Addresses: address,
		Transport: transport,
	})
	if nil != err {
		logrus.Errorf("new elasticsearch client v6 failed. error: %s.", err)
		return
	}
	return
}
