package cache

import (
	"context"

	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type CommandChecker interface {
	Before(client *Client, cmd redis.Cmder)
	After(err error)
}

type CommandCheckerMaker interface {
	NewChecker() CommandChecker
}

type Client struct {
	RedisClient         *redis.Client
	Ctx                 context.Context
	CommandCheckerMaker CommandCheckerMaker
}

func (m *Client) CommandChecker() CommandChecker {
	if m.CommandCheckerMaker == nil {
		return nil
	}

	return m.CommandCheckerMaker.NewChecker()
}

func NewClient(ctx context.Context, redisClient *redis.Client, commCheckerMaker CommandCheckerMaker) (client *Client, err error) {
	client = &Client{
		RedisClient:         redisClient,
		Ctx:                 ctx,
		CommandCheckerMaker: commCheckerMaker,
	}

	err = client.Ping()
	if nil != err {
		logrus.Errorf("ping client failed. %s.", err)
		return
	}

	return
}

func (m *Client) Ping() (err error) {
	_, err = m.RedisClient.Ping().Result()
	if nil != err {
		logrus.Errorf("redis client ping failed. %s.", err)
		return
	}

	return
}

func (m *Client) TTL(key string) (result time.Duration, err error) {
	cmder := m.RedisClient.TTL(key)

	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, cmder)
		defer func() {
			checker.After(err)
		}()
	}

	result, err = cmder.Result()
	return
}

func (m *Client) Expire(key string, expiration time.Duration) (result bool, err error) {
	cmder := m.RedisClient.Expire(key, expiration)

	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, cmder)
		defer func() {
			checker.After(err)
		}()
	}

	result, err = cmder.Result()
	return
}

func (m *Client) Get(key string) (result string, err error) {
	cmder := m.RedisClient.Get(key)

	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, cmder)
		defer func() {
			checker.After(err)
		}()
	}

	result, err = cmder.Result()
	return
}

func (m *Client) Set(key string, value interface{}, expiration time.Duration) (result string, err error) {
	cmder := m.RedisClient.Set(key, value, expiration)

	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, cmder)
		defer func() {
			checker.After(err)
		}()
	}

	result, err = cmder.Result()
	return
}

func (m *Client) Del(keys ...string) (result int64, err error) {
	cmder := m.RedisClient.Del(keys...)

	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, cmder)
		defer func() {
			checker.After(err)
		}()
	}

	result, err = cmder.Result()
	return
}

func (m *Client) HGet(key string, field string) (result string, err error) {
	cmder := m.RedisClient.HGet(key, field)

	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, cmder)
		defer func() {
			checker.After(err)
		}()
	}

	result, err = cmder.Result()
	return
}

func (m *Client) HSet(key string, field string, value interface{}) (result bool, err error) {
	cmder := m.RedisClient.HSet(key, field, value)

	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, cmder)
		defer func() {
			checker.After(err)
		}()
	}

	result, err = cmder.Result()
	return
}
