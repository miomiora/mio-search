package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"mio-search/model"
	"time"
)

func InsertPost(p *model.Post) (err error) {
	err = db.Create(p).Error
	if err != nil {
		zap.L().Error("[dao mysql post] insert post error ", zap.Error(err))
		return
	}
	return nil
}

func InsertPostList(p []*model.Post) (err error) {
	result := db.Create(p)
	err = result.Error
	affected := result.RowsAffected
	fmt.Println("affected: ", affected)
	if err != nil {
		zap.L().Error("[dao mysql post] insert post list error ", zap.Error(err))
		return
	}
	return nil
}

func QueryPostList(page, size int) ([]*model.Post, error) {
	var p []*model.Post
	err := db.Limit(size).Offset(page - 1).Find(&p).Error

	if err != nil {
		return nil, err
	}
	return p, nil
}

func QueryPostListByKeyword(p *model.SearchDTO) ([]*model.Post, error) {
	var post []*model.Post
	err := db.Limit(p.Size).Offset(p.Page-1).Where("content like ? or title like ?", "%"+p.Text+"%", "%"+p.Text+"%").Find(&post).Error

	if err != nil {
		return nil, err
	}
	return post, nil
}

func QueryPostListByIds(p *model.SearchDTO, ids []int64) ([]*model.Post, error) {
	var post []*model.Post
	err := db.Limit(p.Size).Offset(p.Page-1).Where("post_id in ?", ids).Find(&post).Error

	if err != nil {
		return nil, err
	}
	return post, nil
}

func QueryPostListToES() ([]*model.PostESEntity, error) {
	var p []*model.PostESEntity

	err := db.Model(&model.Post{}).Unscoped().Scan(&p).Error
	if err != nil {
		return nil, err
	}

	return p, nil
}

func QueryPostListNearly() ([]*model.PostESEntity, error) {
	var p []*model.PostESEntity
	err := db.Model(&model.Post{}).Where("updated_at > ?", time.Now().Add(time.Minute*-3)).Scan(&p).Error
	if err != nil {
		return nil, err
	}
	return p, nil
}
