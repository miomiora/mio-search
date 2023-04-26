package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mio-search/config"
	"mio-search/model"
	"time"
)

var (
	db *gorm.DB
)

func Init(cfg *config.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		zap.L().Error("[dao mysql] 连接Mysql数据库失败, error = ", zap.Error(err))
		return
	}

	err = db.AutoMigrate(&model.User{}, &model.Post{})
	if err != nil {
		fmt.Println("[dao mysql] 创建表失败！")
	}

	conn, err := db.DB()
	if err != nil {
		zap.L().Error("[dao mysql] 获取sql实例失败！", zap.Error(err))
		return
	}
	conn.SetMaxOpenConns(cfg.MaxOpenConn)
	conn.SetMaxIdleConns(cfg.MaxIdleConn)
	conn.SetConnMaxLifetime(time.Hour * 4)
	return
}

func Close() {
	conn, _ := db.DB()
	_ = conn.Close()
}
