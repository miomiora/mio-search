package datasource

import (
	"mio-search/model"
	"mio-search/util"
)

// iDatasource 新接入数据源必须实现的接口
type iDatasource interface {
	DoSearch(*model.SearchDTO) (interface{}, error)
}

// CenterMap Datasource 注册中心，对外暴露一个单例
var CenterMap = map[string]iDatasource{
	util.Post:    new(postDatasource),
	util.User:    new(userDatasource),
	util.Picture: new(pictureDatasource),
	util.All:     new(allDatasource),
}
