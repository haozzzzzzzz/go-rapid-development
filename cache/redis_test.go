package cache

import (
	"os"
	"testing"

	"context"

	"fmt"

	"log"

	"time"

	"github.com/go-redis/redis"
)

var rawRedisClient *redis.Client
var redisClient *Client

func TestMain(m *testing.M) {
	var err error
	config := &RedisConfigFormat{
		Address: "127.0.0.1:6379",
	}
	rawRedisClient, err = NewRedisClient(config)
	if nil != err {
		log.Fatal(err)
		return
	}

	redisClient = &Client{
		RedisClient: rawRedisClient,
		Ctx:         context.Background(),
		Config:      config,
	}

	os.Exit(m.Run())
}

func TestClient_LRangePop(t *testing.T) {
	key := "test_l_range_pop"
	result, err := redisClient.LRangePop(key, 1)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(result)
}

func TestClient_RPush(t *testing.T) {
	key := "test_l_range_pop"
	result, err := redisClient.RPush(key, "1", "2", "3")
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(result)
}

func TestClient_Set(t *testing.T) {
	key := "test_set"
	result, err := redisClient.Set(key, "1", 3600*time.Second)
	if nil != err {
		t.Error(err)
		return
	}

	fmt.Println(result)
}
