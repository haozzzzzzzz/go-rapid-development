package cache

import (
	"fmt"
	"testing"
)

func TestNewRedisClient(t *testing.T) {
	client, err := NewRedisClient(nil, &RedisConfigFormat{
		Address: "127.0.0.1:6379",
	}, nil, nil)
	if nil != err {
		t.Error(err)
		return
	}

	val, err := client.Set("hello", "world", -1).Result()
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(val)
}
