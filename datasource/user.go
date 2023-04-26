package datasource

import (
	"mio-search/dao/mysql"
	"mio-search/model"
)

type userDatasource struct {
}

func (d userDatasource) DoSearch(p *model.SearchDTO) (interface{}, error) {
	data, err := mysql.QueryUserListByText(p)
	if err != nil {
		return nil, err
	}
	var userList []*model.UserVO
	for _, value := range data {
		user := &model.UserVO{
			UserId:      value.UserId,
			Account:     value.Account,
			Token:       nil,
			Description: value.Description,
			UserRole:    value.UserRole,
		}
		userList = append(userList, user)
	}

	all := &model.AllDatasourceVO{
		Post:    nil,
		Picture: nil,
		User:    userList,
	}
	return all, nil
}
