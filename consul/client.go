package consul

import (
	"encoding/json"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// get
const Nil = uerrors.StringError("consul: nil")

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

// get string
func (m *Client) GetString(key string) (value string, err error) {
	pair, _, err := m.Api.KV().Get(key, nil)
	if nil != err {
		logrus.Errorf("get key value pair failed. key: %s. error: %s.", key, err)
		return
	}

	if pair == nil {
		err = Nil
		return
	}

	value = string(pair.Value)

	return
}

// get yaml
func (m *Client) GetYaml(key string, obj interface{}) (err error) {
	value, err := m.GetString(key)
	if nil != err {
		logrus.Errorf("get consul value failed. error: %s.", err)
		return
	}

	err = yaml.Unmarshal([]byte(value), obj)
	if nil != err {
		logrus.Errorf("unmarshal consul value to yaml obj failed. error: %s.", err)
		return
	}

	return
}

func (m *Client) GetJson(key string, obj interface{}) (err error) {
	value, err := m.GetString(key)
	if nil != err {
		logrus.Errorf("get consul value failed. error: %s.", err)
		return
	}

	err = json.Unmarshal([]byte(value), obj)
	if nil != err {
		logrus.Errorf("unmarshal consul value to json obj failed. error: %s.", err)
		return
	}

	return
}

// put
func (m *Client) PutJson(key string, obj interface{}, sesId string) (err error) {
	value, err := json.Marshal(obj)
	if nil != err {
		logrus.Errorf("marshal obj failed. error: %s.", err)
		return
	}

	meta, err := m.Api.KV().Put(&api.KVPair{
		Key:     key,
		Value:   value,
		Session: sesId,
	}, nil)
	if nil != err {
		logrus.Errorf("put value failed. error: %s.", err)
		return
	}

	_ = meta

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

		logrus.Infof("consul plan ack. %v", string(pair.Value))

		err = localValue.Set(pair.Value)
		if nil != err {
			logrus.Errorf("set local value failed. error: %s.", err)
			return
		}
	}

	go func() {
		logrus.Infof("running consul plan. key: %s, address: %s", key, m.Config.Address)
		err := plan.Run(m.Config.Address)
		if nil != err {
			logrus.Errorf("run consul plan failed. key: %s. address: %s. error: %s.", key, m.Config.Address, err)
			return
		}
	}()

	return
}
