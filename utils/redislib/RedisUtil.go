package redislib

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
	"sync"
)

var redisDB *redis.Client
var once sync.Once

func initRedisDB() (c *redis.Client, err error) {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     viper.GetString(`redis.addr`) + ":" + viper.GetString(`redis.port`),
		Password: viper.GetString(`redis.password`), // no password set
		DB:       viper.GetInt(`redis.db`),          // use default DB
	})

	ctx := context.Background()

	pong, _err := redisDB.Ping(ctx).Result()
	log.Println("redis connect " + pong)
	if _err != nil {
		panic(_err)
	}
	return redisDB, err
}

// GetRedisInstance 获取 redis client的单例
func GetRedisInstance() (c *redis.Client) {
	once.Do(func() {
		redisDB, _ = initRedisDB()
	})
	return redisDB
}

func CloseRedisDB() error {
	return GetRedisInstance().Close()
}
