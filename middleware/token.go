package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mio-search/controller"
	"mio-search/dao/redis"
)

func RefreshToken(c *gin.Context) {
	// 从请求头中获取Token, 没有token就直接返回
	token := c.GetHeader("Authorization")
	if token == "" {
		c.Next()
	}
	userId, err := redis.QueryTokenExist(token)
	if err == nil {
		// 走到这步意味着 token 有效，并且已经被刷新
		c.Set(redis.KeyUserId, userId)
	}
}

func AuthToken(c *gin.Context) {
	// 直接验证是否保存了 user_id 的字段
	// 获取userId
	value, exist := c.Get(redis.KeyUserId)
	if !exist {
		zap.L().Error("[middleware token] get userId error ")
		controller.ResponseError(c, controller.ErrorNotLogin)
		c.Abort()
		return
	}
	_, ok := value.(int64)
	if !ok {
		zap.L().Error("[middleware token] get userId error ")
		controller.ResponseError(c, controller.ErrorServerBusy)
		c.Abort()
		return
	}
	c.Next()
}
