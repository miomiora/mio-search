package logic

import (
	"mio-search/dao/mysql"
	"mio-search/model"
	"mio-search/util"
)

func InsertPost(p *model.PostDTOInsert, userId int64) (err error) {
	postId := util.GenSnowflakeID()
	// 构造一个User实例
	post := &model.Post{
		UserId:  userId,
		PostId:  postId,
		Content: p.Content,
		Title:   p.Title,
	}
	// 保存进数据库
	err = mysql.InsertPost(post)
	return
}

func GetPostList(p *model.SearchDTO) ([]*model.PostVO, error) {
	data, err := mysql.QueryPostListByKeyword(p)
	if err != nil {
		return nil, err
	}

	var postList []*model.PostVO
	var user *model.User

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
	return postList, nil
}
