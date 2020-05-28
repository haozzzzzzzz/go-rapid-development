package consul

import (
	"encoding/json"
	"fmt"
	"github.com/haozzzzzzzz/go-rapid-development/utils/uerrors"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/sirupsen/logrus"
	"github.com/xiam/to"
	"gopkg.in/yaml.v2"
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

type WatchServiceCallback func(heathChecks []*api.HealthCheck) (err error)

// watch service
func (m *Client) WatchChecks(serviceName string, callback WatchServiceCallback) (err error) {
	plan, err := watch.Parse(map[string]interface{}{
		"type":    "checks",
		"service": serviceName,
	})
	if nil != err {
		logrus.Errorf("parse watch checks params error: %s", err)
		return
	}

	plan.HybridHandler = func(blockVal watch.BlockingParamVal, val interface{}) {
		healthChecks, ok := val.([]*api.HealthCheck)
		if !ok {
			return
		}

		errCallback := callback(healthChecks)
		if errCallback != nil {
			logrus.Errorf("watch checks callback error: %s", err)
			return
		}
	}

	go func() {
		errRun := plan.Run(m.Config.Address)
		if nil != errRun {
			logrus.Errorf("error: %s", err)
		}
	}()
	return
}

// 注册服务
func (m *Client) RegisterService(
	serviceName string,
	ip string, // ip加端口
	port string,
	metricsPath string,
	checkInterval string,
	tags []string,
	meta map[string]string,
) (err error) {
	address := fmt.Sprintf("%s:%s", ip, port)
	serviceId := address
	metricsAddress := fmt.Sprintf("%s%s", address, metricsPath)

	if meta == nil {
		meta = make(map[string]string)
	}

	meta["metrics_path"] = metricsPath

	// https://www.consul.io/api/agent/service.html#register-service
	err = m.Api.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name:    serviceName,
		ID:      serviceId,
		Tags:    tags,
		Address: ip,
		Meta:    meta,
		Port:    to.Int(port),

		// https://www.consul.io/api/agent/check.html#register-check
		Checks: api.AgentServiceChecks{
			&api.AgentServiceCheck{
				Name:                           "metrics_check",
				CheckID:                        metricsAddress,
				Interval:                       checkInterval,
				HTTP:                           fmt.Sprintf("http://%s", metricsAddress),
				Status:                         "passing",
				Timeout:                        "10s",
				DeregisterCriticalServiceAfter: "10m",
			},
		},
	})
	if nil != err {
		logrus.Errorf("register service error: %s", err)
	}

	return
}

// 取消注册服务
func (m *Client) DeregisterService(
	ip string,
	port string,
) (err error) {
	address := fmt.Sprintf("%s:%s", ip, port)
	serviceId := address

	err = m.Api.Agent().ServiceDeregister(serviceId)
	if nil != err {
		logrus.Errorf("deregister service error: %s", err)
	}

	return
}
