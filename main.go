package main

import (
	"context"
	"fmt"
	"mio-search/config"
	"mio-search/dao/es"
	"mio-search/dao/mysql"
	"mio-search/dao/redis"
	"mio-search/logger"
	"mio-search/router"
	"mio-search/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// Go Web 开发通用脚手架模板

func main() {
	// 加载配置
	if err := config.Init(); err != nil {
		fmt.Printf("init config error : %s \n", err)
		return
	}

	// 初始化日志
	if err := logger.Init(config.Conf.LogConfig, config.Conf.Mode); err != nil {
		fmt.Printf("init logger error : %s \n", err)
		return
	}
	defer zap.L().Sync()

	// 初始化 MySQL
	if err := mysql.Init(config.Conf.MySQLConfig); err != nil {
		fmt.Printf("init MySQL error : %s \n", err)
		return
	}
	defer mysql.Close()

	// 初始化 Redis
	if err := redis.Init(config.Conf.RedisConfig); err != nil {
		fmt.Printf("init Redis error  %s \n", err)
		return
	}
	defer redis.Close()

	// 初始化雪花算法
	if err := util.Init(config.Conf.StartTime, config.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake error  %s \n", err)
		return
	}

	// 初始化 post 的数据
	//if err := spider.FetchInitPostList(); err != nil {
	//	fmt.Printf("fetch post data spider error  %s \n", err)
	//	return
	//}

	// 初始化 ElasticSearch
	if err := es.InitElasticSearch(config.Conf.ElasticSearchConfig); err != nil {
		fmt.Printf("init ElasticSearch error  %s \n", err)
		return
	}

	// 执行同步 MySQL 与 ElasticSearch 的数据
	go es.SyncPostListES()

	// 注册路由
	r := router.Setup()
	// 启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Conf.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Info("listen: %s\n", zap.Error(err))
		}
	}()

	// 为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5 秒执行最后的请求
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
