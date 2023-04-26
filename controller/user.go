package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mio-search/logic"
	"mio-search/model"
)

func UserLogin(c *gin.Context) {
	// 1、校验参数
	u := new(model.UserDTOLogin)
	if err := c.ShouldBindJSON(u); err != nil {
		// 请求参数有误
		zap.L().Error("[controller user] login with invalid param ", zap.Error(err))
		ResponseError(c, ErrorInvalidParams)
		return
	}
	// 2、业务处理
	data, err := logic.UserLogin(u)
	if err != nil {
		zap.L().Error("[controller user] login failed ", zap.String("Account", u.Account), zap.Error(err))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ResponseErrorWithMsg(c, ErrorInvalidParams, "用户名或密码错误！")
			return
		}
		ResponseError(c, ErrorServerBusy)
		return
	}
	// 3、返回响应
	ResponseOK(c, data)
}

func UserRegister(c *gin.Context) {
	// 1、参数校验
	u := new(model.UserDTORegister)
	if err := c.ShouldBindJSON(u); err != nil {
		// 请求参数有误
		zap.L().Error("[controller user] register with invalid param ", zap.Error(err))
		ResponseError(c, ErrorInvalidParams)
		return
	}
	// 2、业务处理
	if err := logic.UserRegister(u); err != nil {
		zap.L().Error("register failed ", zap.Error(err))
		if errors.Is(err, logic.ErrorUserExist) {
			ResponseErrorWithMsg(c, ErrorInvalidParams, err.Error())
			return
		}
		ResponseError(c, ErrorServerBusy)
		return
	}
	// 3、返回响应
	ResponseOK(c, nil)
}

func GetUserList(c *gin.Context) {
	// 参数校验
	// 1、获取参数，绑定query参数到结构体上
	// 初始化结构体指定初始参数
	p := new(model.SearchDTO)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("[controller post] get user list with invalid query params", zap.Error(err))
	}
	// 2、业务处理
	data, err := logic.GetUserList(p)
	if err != nil {
		zap.L().Error("logic get user list error ", zap.Error(err))
		ResponseError(c, ErrorServerBusy)
		return
	}
	ResponseOK(c, data)
}
