package consul

import (
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
	"github.com/sirupsen/logrus"
)

type LocalValue interface {
	// should add lock, multi goroutines access
	Set(value []byte) error
}

type Client struct {
	Api    *api.Client
	Config *ClientConfigFormat
}

func NewClient(config *ClientConfigFormat) (client *Client, err error) {
	client = &Client{
		Config: config,
	}
	client.Api, err = api.NewClient(&api.Config{
		Address: config.Address,
	})
	if nil != err {
		logrus.Errorf("new consul api client failed. error: %s.", err)
		return
	}

	return
}

// sync get
func (m *Client) GetSync(key string, localValue LocalValue) (err error) {
	pair, _, err := m.Api.KV().Get(key, nil)
	if nil != err {
		logrus.Errorf("get key value pair failed. error: %s.", err)
		return
	}

	if pair == nil {
		err = uerrors.Newf("consul key not exist. %s", key)
		return
	}

	err = localValue.Set(pair.Value)
	if nil != err {
		logrus.Errorf("set local value failed. error: %s.", err)
		return
	}

	return
}

// async watch
func (m *Client) Watch(key string, localValue LocalValue) (err error) {
	plan, err := watch.Parse(map[string]interface{}{
		"type": "key",
		"key":  key,
	})
	if nil != err {
		logrus.Errorf("parse consul watch plan failed. error: %s.", err)
		return
	}

	plan.Handler = func(u uint64, i interface{}) {
		if i == nil {
			logrus.Warnf("plan return value is nil. key: %s", key)
			return
		}

		pair, ok := i.(*api.KVPair)
		if !ok {
			logrus.Warnf("plan return value's type is not KVPair. %v", i)
			return
		}

		err = localValue.Set(pair.Value)
		if nil != err {
			logrus.Errorf("set local value failed. error: %s.", err)
			return
		}
	}

	go plan.Run(m.Config.Address)
	return
}
