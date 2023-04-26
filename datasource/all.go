package datasource

import (
	"mio-search/model"
	"mio-search/util"
)

type allDatasource struct {
}

func (a allDatasource) DoSearch(p *model.SearchDTO) (interface{}, error) {

	postList, err := CenterMap[util.Post].DoSearch(p)
	if err != nil {
		return nil, err
	}
	userList, err := CenterMap[util.User].DoSearch(p)
	if err != nil {
		return nil, err
	}
	pictureList, err := CenterMap[util.Picture].DoSearch(p)
	if err != nil {
		return nil, err
	}

	all := &model.AllDatasourceVO{
		Post:    postList,
		Picture: pictureList,
		User:    userList,
	}

	return all, nil
}
