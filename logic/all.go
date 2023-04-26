package logic

import (
	"mio-search/datasource"
	"mio-search/model"
)

func SearchAll(p *model.SearchDTO) (interface{}, error) {
	return datasource.CenterMap[p.Type].DoSearch(p)
}
