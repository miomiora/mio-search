package router

import (
	"github.com/gin-gonic/gin"
	"mio-search/controller"
	"mio-search/logger"
	"mio-search/middleware"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middleware.Cors, middleware.RefreshToken)
	v1 := r.Group("/api")

	// 登录注册
	v1.POST("/login", controller.UserLogin)
	v1.POST("/register", controller.UserRegister)

	// 以下操作需要登录
	v1.Use(middleware.AuthToken)
	v1.GET("/user", controller.GetUserList)

	v1.POST("/post", controller.InsertPost)
	v1.GET("/post", controller.GetPostList)

	v1.GET("/picture", controller.GetPictureList)

	v1.POST("/all", controller.SearchAll)

	r.NoRoute(func(c *gin.Context) {
		controller.ResponseError(c, controller.ErrorNotFound)
	})
	return r
}
