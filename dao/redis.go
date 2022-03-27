package dao

import (
	"HomeWorkGo/setting"
	"github.com/go-redis/redis"
)

var (
	RDB *redis.Client
)

func InitRedis(cfg *setting.RedisConfig) (err error) {

	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password, // no password set
		DB:       0,            // use default DB
	})
	_, err = RDB.Ping().Result()

	return err

}
