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

// wait lock
func (m *RedisLock) Lock(
	key string,
	wait bool,
	waitDur time.Duration,
	expire time.Duration,
) (err error) {
	key = m.KeyLock(key)
	var ok bool

	if expire == 0 {
		expire = TTL_REDIS_LOCK
	}

	for i := 0; i < 100; i++ {
		ok, err = m.Client.SetNX(key, "LOCKED", expire) // 尝试上锁
		if nil != err {
			logrus.Errorf("set redis lock failed. %s.", err)
			break
		}

		if ok || false == wait {
			break
		}

		if wait && waitDur > 0 {
			time.Sleep(waitDur)
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

// not wait
func (m *RedisLock) LockNotWaitLock(key string, expire time.Duration) (success bool) {
	if expire == 0 {
		expire = TTL_REDIS_LOCK
	}

	success, err := m.Client.SetNX(key, "LOCKED", expire)
	if nil != err {
		logrus.Errorf("set not wait lock failed. error: %s.", err)
		return
	}
	return
}

func (m *RedisLock) UnlockNotWaitLock(key string) {
	_, err := m.Client.Del(key)
	if nil != err {
		logrus.Errorf("unlock not wait lock failed. error: %s.", err)
		return
	}
	return
}
