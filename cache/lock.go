package cache

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type RedisLock struct {
	Client *Client
}

func NewRedisLock(client *Client) *RedisLock {
	return &RedisLock{
		Client: client,
	}
}

const TTL_REDIS_LOCK = 60 * time.Second

func (m *RedisLock) KeyLock(key string) string {
	return fmt.Sprintf("LOCK:%s", key)
}

func (m *RedisLock) Lock(key string, wait bool) (err error) {
	key = m.KeyLock(key)
	var ok bool
	for i := 0; i < 100; i++ {
		ok, err = m.Client.SetNX(key, "LOCKED", TTL_REDIS_LOCK) // 尝试上锁
		if nil != err {
			logrus.Errorf("set redis lock failed. %s.", err)
			break
		}

		if ok || false == wait {
			break
		}

	}

	if !ok {
		err = errors.New("add redis lock failed")
	}

	return
}

func (m *RedisLock) Unlock(key string) (err error) {
	key = m.KeyLock(key)
	_, err = m.Client.Del(key)
	if nil != err {
		logrus.Errorf("del redis lock failed. %s.", err)
		return
	}
	return
}
