package consul

import (
	"fmt"
	"github.com/haozzzzzzzz/go-rapid-development/utils/yaml"
	"github.com/sirupsen/logrus"
	"testing"

	"time"

	"github.com/hashicorp/consul/api"
)

func TestClient_PutJson(t *testing.T) {
	client, err := NewClient(&ClientConfigFormat{
		Address: "127.0.0.1:8500",
	})
	if nil != err {
		t.Error(err)
		return
	}

	err = client.PutJson("dev/test/helloworld", &struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{
		Key:   "hello",
		Value: "world",
	}, "")
	if nil != err {
		t.Error(err)
		return
	}

}

func TestClient_GetJson(t *testing.T) {
	client, err := NewClient(&ClientConfigFormat{
		Address: "127.0.0.1:8500",
	})
	if nil != err {
		t.Error(err)
		return
	}

	m := make(map[string]interface{})
	err = client.GetJson("dev/test/helloworld1", &m)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(m)
}

func TestAcquire(t *testing.T) {
	client, err := NewClient(&ClientConfigFormat{
		Address: "127.0.0.1:8500",
	})
	if nil != err {
		t.Error(err)
		return
	}

	sesId, _, err := client.Api.Session().Create(&api.SessionEntry{
		LockDelay: 15 * time.Second,
		TTL:       "30s",
	}, nil)
	if nil != err {
		t.Error(err)
		return
	}

	acquire, _, err := client.Api.KV().Acquire(&api.KVPair{
		Key:     "dev/lock/test_lock",
		Value:   []byte("1"),
		Session: sesId,
	}, nil)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(acquire)

	//_, err = client.Api.KV().Delete("dev/lock/test_lock", nil)
	//if nil != err {
	//	t.Error(err)
	//	return
	//}

	release, _, err := client.Api.KV().Release(&api.KVPair{
		Key:     "dev/lock/test_lock",
		Value:   []byte("0"),
		Session: sesId,
	}, nil)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(release)
}

func TestClient_GetYaml(t *testing.T) {
	type LogOutputConfigFormat struct {
		Filedir    string `json:"filedir" yaml:"filedir" validate:"required"`
		MaxSize    int    `json:"max_size" yaml:"max_size" validate:"required"`
		MaxBackups int    `json:"max_backups" yaml:"max_backups" validate:"required"`
		MaxAge     int    `json:"max_age" yaml:"max_age" validate:"required"`
		Compress   bool   `json:"compress" yaml:"compress"`
	}

	type LogConfigFormat struct {
		LogLevel logrus.Level           `json:"log_level" yaml:"log_level" validate:"required"`
		Output   *LogOutputConfigFormat `json:"output" yaml:"output"`
	}

	client, err := NewClient(&ClientConfigFormat{
		Address: "127.0.0.1:8500",
	})
	if nil != err {
		t.Error(err)
		return
	}
	_ = client

	config := &LogConfigFormat{}
	//err = client.GetYaml("dev/service/video_buddy_user_device/log.yaml", config)
	//if nil != err {
	//	t.Error(err)
	//	return
	//}
	err = yaml.ReadYamlFromFile("/Users/hao/Documents/Projects/XunLei/video_buddy_user_device/common/stage/dev/config/consul/log.yaml", config)
	if nil != err {
		t.Error(err)
		return
	}
}