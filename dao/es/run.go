package es

import (
	"go.uber.org/zap"
	"mio-search/dao/mysql"
	"time"
)

// SyncPostListES 同步 MySQL 的数据到 ES 中
func SyncPostListES() {
	go fullUpdate()
	go addUpdate()
}

// 全量更新只执行一次
func fullUpdate() {
	for {
		postList, err := mysql.QueryPostListToES()
		if err != nil {
			continue
		}
		err = InsertPost(postList)
		if err != nil {
			continue
		}
		break
	}
}

// 增量更新 1 分钟执行一次
func addUpdate() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			func() {
				nearly, err2 := mysql.QueryPostListNearly()
				if err2 != nil {
					zap.L().Info("[dao es post error] query post list nearly error ", zap.Error(err2))
					return
				}
				if nearly == nil || len(nearly) == 0 {
					return
				}
				err2 = InsertPost(nearly)
				if err2 != nil {
					zap.L().Info("[dao es post error] insert post list to elastic search ", zap.Error(err2))
					return
				}
			}()
		}
	}
}
