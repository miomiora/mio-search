package redis

import (
	"go.uber.org/zap"
	"strconv"
)

func InsertTokenByUserId(token string, userId int64) (err error) {
	// 使用 pipeline 减少 RTT
	pipeline := client.TxPipeline()

	// 把 token 插入到 redis中
	key := TokenPrefix + token
	pipeline.HSet(ctx, key, KeyUserId, userId)
	// 为 token 设置过期时间
	pipeline.Expire(ctx, key, TokenTimeout)

	// 执行 pipeline
	_, err = pipeline.Exec(ctx)

	return
}

func QueryTokenExist(token string) (int64, error) {
	key := TokenPrefix + token

	id, err := client.HGet(ctx, key, KeyUserId).Result()
	if err != nil {
		zap.L().Error("[middleware token] client hget key ", zap.Error(err))
		return -1, err
	}

	err = client.Expire(ctx, key, TokenTimeout).Err()
	if err != nil {
		zap.L().Error("[middleware token] client expire key ", zap.Error(err))
		return -1, err
	}

	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		zap.L().Error("[middleware token] strconv ParseInt id error ", zap.Error(err))
		return -1, err
	}

	return i, nil
}
