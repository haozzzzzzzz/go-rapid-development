package cache

import (
	"context"

	"time"

	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
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

func (m *Client) Keys(pattern string) (result []string, err error) {
	cmder := m.RedisClient.Keys(pattern)
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

func (m *Client) MGet(keys ...string) (result []interface{}, err error) {
	cmder := m.RedisClient.MGet(keys...)

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

func (m *Client) MSet(pairs ...interface{}) (result string, err error) {
	cmder := m.RedisClient.MSet(pairs...)

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

// empty string for nil value
func (m *Client) MGetString(key string) (result []string, err error) {
	iResults, err := m.MGet(key)

	for _, iResult := range iResults {
		var strResult string

		if iResult != nil {
			strResult = iResult.(string)
		}

		result = append(result, strResult)
	}

	return
}

func (m *Client) GetJson(key string, obj interface{}) (err error) {
	result, err := m.Get(key)
	if nil != err {
		if err != redis.Nil {
			logrus.Errorf("redis get %s failed. %s.", key, err)
		}
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

// incr
func (m *Client) Incr(key string) (result int64, err error) {
	cmder := m.RedisClient.Incr(key)

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

func (m *Client) IncrBy(key string, value int64) (result int64, err error) {
	cmder := m.RedisClient.IncrBy(key, value)

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

func (m *Client) HMGet(key string, fields ...string) (result []interface{}, err error) {
	cmder := m.RedisClient.HMGet(key, fields...)

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

func (m *Client) HMSet(key string, fields map[string]interface{}) (result string, err error) {
	cmder := m.RedisClient.HMSet(key, fields)

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

func (m *Client) HIncrBy(key string, field string, incr int64) (result int64, err error) {
	cmder := m.RedisClient.HIncrBy(key, field, incr)
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

func (m *Client) HGetAll(key string) (result map[string]string, err error) {
	cmder := m.RedisClient.HGetAll(key)
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, cmder)
		defer func() {
			checker.After(err)
		}()
	}

	result, err = cmder.Result()
	if nil != err {
		logrus.Errorf("hgetall failed. error: %s.", err)
		return
	}

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

// 没有则返回redis.Nil
func (m *Client) HGetJSON(key string, field string, value interface{}) (err error) {
	result, err := m.HGet(key, field)
	if nil != err {
		if err != redis.Nil {
			logrus.Errorf("get hash item failed. key: %s. field: %s. error: %s.", key, field, err)
		}
		return
	}

	err = json.Unmarshal([]byte(result), value)
	if nil != err {
		logrus.Errorf("unmarshal hash item failed. error: %s.", err)
		return
	}

	return
}

// set
func (m *Client) SAdd(key string, members ...interface{}) (result int64, err error) {
	cmder := m.RedisClient.SAdd(key, members...)

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

func (m *Client) SRem(key string, members ...interface{}) (result int64, err error) {
	cmder := m.RedisClient.SRem(key, members...)

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

func (m *Client) SMembers(key string) (result []string, err error) {
	cmder := m.RedisClient.SMembers(key)
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

// key不存在，不会返回redis.Nil
func (m *Client) SRandMemberN(key string, count int64) (result []string, err error) {
	cmder := m.RedisClient.SRandMemberN(key, count)

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

// sorted set
func (m *Client) ZCard(key string) (result int64, err error) {
	cmder := m.RedisClient.ZCard(key)

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

func (m *Client) ZUnionStore(destination string, store redis.ZStore, keys ...string) (result int64, err error) {
	cmder := m.RedisClient.ZUnionStore(destination, store, keys...)
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

func (m *Client) ZInterStore(destination string, store redis.ZStore, keys ...string) (result int64, err error) {
	cmder := m.RedisClient.ZInterStore(destination, store, keys...)
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

// list
func (m *Client) RPush(key string, values ...interface{}) (result int64, err error) {
	cmder := m.RedisClient.RPush(key, values...)
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

func (m *Client) LPop(key string) (result string, err error) {
	cmder := m.RedisClient.LPop(key)
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

// 非原子操作的批量将队列头部的元素弹出
func (m *Client) LRangePop(key string, num uint64) (result []string, err error) {
	if num <= 0 {
		err = errors.New("LRangePop num params should be greater than 0")
		return
	}

	result = make([]string, 0)
	defer func() {
		if nil != err {
			result = make([]string, 0)
			return
		}
	}()

	stop := num - 1
	result, err = m.LRange(key, int64(0), int64(stop))
	if nil != err {
		logrus.Errorf("range get list failed. key: %s, num: %d. error: %s.", key, num, err)
		return
	}

	if len(result) == 0 {
		return
	}

	_, err = m.LTrim(key, int64(stop+1), int64(-1))
	if nil != err {
		logrus.Errorf("trim list failed. %s.", err)
		return
	}

	return
}

func (m *Client) LLen(key string) (result int64, err error) {
	cmder := m.RedisClient.LLen(key)
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

func (m *Client) LRange(key string, start int64, stop int64) (result []string, err error) {
	cmder := m.RedisClient.LRange(key, start, stop)
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

func (m *Client) LTrim(key string, start int64, stop int64) (result string, err error) {
	cmder := m.RedisClient.LTrim(key, start, stop)
	checker := m.CommandChecker()
	if checker != nil {
		checker.Before(m, cmder)
		defer func() {
			checker.After(err)
		}()
	}

	result, err = cmder.Result()

	m.RedisClient.Pipeline()
	return
}
