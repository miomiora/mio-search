package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"mio-search/config"
	"time"
)

var (
	client *redis.Client
	ctx    context.Context
)

const (
	TokenPrefix  = "login:token:"
	TokenTimeout = time.Hour * 24 * 7

	KeyUserId = "user_id"
)

func Init(cfg *config.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // 密码
		DB:       cfg.Db,       // 数据库
		PoolSize: cfg.PoolSize, // 连接池大小
	})

	_, err = client.Ping(context.Background()).Result()
	ctx = context.Background()
	return
}

func Close() {
	_ = client.Close()
}
