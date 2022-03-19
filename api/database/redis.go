package database

import (
	"apartment/config"
	"sync"

	"github.com/go-redis/redis"
)

type Redis struct {
	*redis.Client
}

var redisInstance *Redis
var once sync.Once

func GetRedis() *Redis {
	once.Do(func() {
		redisInstance = &Redis{
			Client: redis.NewClient(&redis.Options{
				Addr:     config.Radis.Address,
				Password: config.Radis.Password,
				DB:       config.Radis.DB,
			}),
		}
	})
	return redisInstance
}
