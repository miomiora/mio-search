package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mio-search/logic"
	"mio-search/model"
)

func GetPictureList(c *gin.Context) {
	// 参数校验
	// 1、获取参数，绑定query参数到结构体上
	// 初始化结构体指定初始参数
	p := new(model.SearchDTO)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("[controller post] get user list with invalid query params", zap.Error(err))
	}
	// 2、业务处理
	data, err := logic.GetPictureList(p)
	if err != nil {
		zap.L().Error("logic get post list error ", zap.Error(err))
		ResponseError(c, ErrorServerBusy)
		return
	}
	ResponseOK(c, data)
}
