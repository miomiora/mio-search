package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mio-search/dao/redis"
	"mio-search/logic"
	"mio-search/model"
)

func InsertPost(c *gin.Context) {
	// 1、参数校验
	p := new(model.PostDTOInsert)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误
		zap.L().Error("[controller user] insert post with invalid param ", zap.Error(err))
		ResponseError(c, ErrorInvalidParams)
		return
	}
	// 获取userId
	value, exist := c.Get(redis.KeyUserId)
	if !exist {
		zap.L().Error("[controller post] get userId error ")
		ResponseError(c, ErrorServerBusy)
		return
	}
	userId, ok := value.(int64)
	if !ok {
		zap.L().Error("[controller post] get userId error ")
		ResponseError(c, ErrorServerBusy)
		return
	}
	// 2、业务处理
	if err := logic.InsertPost(p, userId); err != nil {
		zap.L().Error("insert post failed ", zap.Error(err))
		ResponseError(c, ErrorServerBusy)
		return
	}
	// 3、返回响应
	ResponseOK(c, nil)
}

func GetPostList(c *gin.Context) {
	// 参数校验
	// 1、获取参数，绑定query参数到结构体上
	// 初始化结构体指定初始参数
	p := new(model.SearchDTO)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("[controller post] get user list with invalid query params", zap.Error(err))
	}
	// 2、业务处理
	data, err := logic.GetPostList(p)
	if err != nil {
		zap.L().Error("logic get post list error ", zap.Error(err))
		ResponseError(c, ErrorServerBusy)
		return
	}
	ResponseOK(c, data)
}
