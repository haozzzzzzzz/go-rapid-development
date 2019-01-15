package consul

import (
	"fmt"
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
