package cache

import (
	"context"

	"time"

	"encoding/json"

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
	Config              *RedisConfigFormat
}

func (m *Client) CommandChecker() CommandChecker {
	if m.CommandCheckerMaker == nil {
		return nil
	}

	return m.CommandCheckerMaker.NewChecker()
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

func (m *Client) GetJson(key string, obj interface{}) (err error) {
	result, err := m.Get(key)
	if nil != err {
		logrus.Errorf("redis get %s failed. %s.", key, err)
		return
	}

	err = json.Unmarshal([]byte(result), obj)
	if nil != err {
		logrus.Errorf("json unmarshal redis value to obj failed. %s.", err)
		return
	}

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

func (m *Client) SetNX(key string, value interface{}, expiration time.Duration) (result bool, err error) {
	cmder := m.RedisClient.SetNX(key, value, expiration)

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

func (m *Client) HDel(key string, fields ...string) (result int64, err error) {
	cmder := m.RedisClient.HDel(key, fields...)
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

func (m *Client) HKeys(key string) (result []string, err error) {
	cmder := m.RedisClient.HKeys(key)
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

func (m *Client) SetJSON(key string, value interface{}, expiration time.Duration) (result string, err error) {
	byteValue, err := json.Marshal(value)
	if nil != err {
		logrus.Errorf("json marshal failed. %s.", err)
		return
	}

	result, err = m.Set(key, string(byteValue), expiration)
	if nil != err {
		logrus.Errorf("redis set failed. %s.", err)
		return
	}

	return
}

func (m *Client) HSetJSON(key string, field string, value interface{}) (result bool, err error) {
	byteValue, err := json.Marshal(value)
	if nil != err {
		logrus.Errorf("json marshal failed. %s.", err)
		return
	}

	result, err = m.HSet(key, field, string(byteValue))
	if nil != err {
		logrus.Errorf("redis hset failed. %s.", err)
		return
	}
	return
}

// sorted set
func (m *Client) ZAdd(key string, members ...redis.Z) (result int64, err error) {
	cmder := m.RedisClient.ZAdd(key, members...)
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

func (m *Client) ZRem(key string, members ...interface{}) (result int64, err error) {
	cmder := m.RedisClient.ZRem(key, members...)
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

func (m *Client) ZRevRange(key string, start int64, stop int64) (result []string, err error) {
	cmder := m.RedisClient.ZRevRange(key, start, stop)
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

func (m *Client) ZUnionStore(dest string, store redis.ZStore, keys ...string) (result int64, err error) {
	cmder := m.RedisClient.ZUnionStore(dest, store, keys...)
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
