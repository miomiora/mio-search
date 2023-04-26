package datasource

import (
	"mio-search/model"
	"mio-search/spider"
)

type pictureDatasource struct {
}

func (d pictureDatasource) DoSearch(p *model.SearchDTO) (interface{}, error) {
	data, err := spider.GetPicture(p)
	all := &model.AllDatasourceVO{
		Post:    nil,
		Picture: data,
		User:    nil,
	}
	return all, err
}
