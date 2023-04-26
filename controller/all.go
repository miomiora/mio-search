package controller

import (
	"github.com/gin-gonic/gin"
	"mio-search/logic"
	"mio-search/model"
)

func SearchAll(c *gin.Context) {
	// 1、校验参数
	searchParams := new(model.SearchDTO)
	err := c.ShouldBindJSON(searchParams)
	if err != nil {
		ResponseError(c, ErrorInvalidParams)
		return
	}
	// 2、业务处理
	data, err := logic.SearchAll(searchParams)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorServerBusy, err.Error())
		return
	}
	// 3、返回响应
	ResponseOK(c, data)
}
