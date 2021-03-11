package cache

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type RedisConfigFormat struct {
	Address  string `json:"address" yaml:"address" validate:"required"`
	Password string `json:"password" yaml:"password"`
	PoolSize int    `json:"pool_size" yaml:"pool_size"` // 默认是CPU核心数*10
	DB       int    `json:"db" yaml:"db"`
}

func NewRedisClient(config *RedisConfigFormat) (client *redis.Client, err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: config.Password,
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

// redis://<user>:<pass>@localhost:6379/<db>
func NewRedisClientByUrl(redisUrl string) (client *redis.Client, err error) {
	options, err := redis.ParseURL(redisUrl)
	if err != nil {
		logrus.Errorf("parse redis url failed. error: %s", err)
		return
	}

	client = redis.NewClient(options)
	result, err := client.Ping().Result()
	if nil != err {
		logrus.Errorf("ping redis service failed. result: %s, err:%s.", result, err)
		return
	}

	return
}
