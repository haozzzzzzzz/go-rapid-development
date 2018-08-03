package cache

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type RedisConfigFormat struct {
	Address  string `json:"address" binding:"required"`
	PoolSize int    `json:"pool_size" yaml:"pool_size"` // 默认是CPU核心数*10
	DB       int    `json:"db" yaml:"db"`
}

func NewRedisClient(config *RedisConfigFormat) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     config.Address,
		PoolSize: config.PoolSize,
		DB:       config.DB,
	})

	result, err := client.Ping().Result()
	if nil != err {
		logrus.Errorf("ping redis service failed. result: %s, err:%s.", result, err)
		return
	}

	return
}
