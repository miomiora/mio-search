package datasource

import (
	"go.uber.org/zap"
	"mio-search/dao/es"
	"mio-search/dao/mysql"
	"mio-search/model"
)

type postDatasource struct {
}

func (d postDatasource) DoSearch(p *model.SearchDTO) (interface{}, error) {
	postList, err := es.QueryPostByKeyword(p)
	if err != nil {
		zap.L().Info("[datasource post do search error] elastic search error ", zap.Error(err))
		// 调用原生 MySQL 搜索
		return d.doSearchByMySQL(p)
	}

	var ids []int64
	for _, entity := range postList {
		ids = append(ids, entity.PostId)
	}
	data, err := mysql.QueryPostListByIds(p, ids)
	if err != nil {
		return nil, err
	}

	for _, post := range data {
		for _, highlight := range postList {
			if post.PostId == highlight.PostId {
				if highlight.Title == "" {
					continue
				}
				post.Title = highlight.Title
				post.Content = highlight.Content
				continue
			}
		}
	}

	return d.bindUserToPost(data)
}

// 当 ES 搜索失败时调用原生的 MySQL 搜索
func (d postDatasource) doSearchByMySQL(p *model.SearchDTO) (interface{}, error) {

	data, err := mysql.QueryPostListByKeyword(p)
	if err != nil {
		return nil, err
	}

	return d.bindUserToPost(data)
}

func (d postDatasource) bindUserToPost(data []*model.Post) (interface{}, error) {
	var postList []*model.PostVO
	var user *model.User
	var err error

	for _, value := range data {
		user, err = mysql.QueryUserByUserId(value.UserId)
		if err != nil {
			return nil, err
		}

		post := &model.PostVO{
			Account:   user.Account,
			Title:     value.Title,
			Content:   value.Content,
			PostId:    value.PostId,
			CreatedAt: value.CreatedAt,
		}

		postList = append(postList, post)
	}

	all := &model.AllDatasourceVO{
		Post:    postList,
		Picture: nil,
		User:    nil,
	}

	return all, nil
}
