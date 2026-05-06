package db

import "github.com/redis/go-redis/v9"

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

var rdb *redis.Client

func NewRedisClient(cfg *RedisConfig) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password, // 没有密码，默认值
		DB:       cfg.DB,       // 默认DB 0
	})
}

func GetRedisClient() *redis.Client {
	return rdb
}
