package logic

import (
	"mio-search/model"
	"mio-search/spider"
)

func GetPictureList(p *model.SearchDTO) ([]*model.PictureVO, error) {
	return spider.GetPicture(p)
}
