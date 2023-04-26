package spider

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"mio-search/model"
	"net/url"
)

var ErrorOutRange = errors.New("图片不存在")

func GetPicture(p *model.SearchDTO) ([]*model.PictureVO, error) {
	var escape, urlStr string
	if p.Text == "" {
		escape = url.QueryEscape("壁纸")
	} else {
		escape = url.QueryEscape(p.Text)
	}
	urlStr = "https://cn.bing.com/images/search?q=" + escape + "&first=1"
	c := colly.NewCollector()

	// 绑定 字段 m 中的 json 数据
	var picMap map[string]string

	var bingPicture []*model.PictureVO

	c.OnHTML(".iuscp.isv", func(e *colly.HTMLElement) {
		pic := e.ChildAttr(".iusc", "m")
		title := e.ChildAttr(".inflnk", "aria-label")

		err := json.Unmarshal([]byte(pic), &picMap)
		if err != nil {
			fmt.Println("Unmarshal error ", err)
			return
		}

		bingPicture = append(bingPicture, &model.PictureVO{
			Title:   title,
			Picture: picMap["turl"],
			Purl:    picMap["purl"],
		})
	})

	err := c.Visit(urlStr)
	if err != nil {
		fmt.Println("c.Visit error ", err)
		return nil, err
	}

	start := (p.Page - 1) * p.Size
	end := start + p.Size

	if end >= len(bingPicture) {
		return nil, ErrorOutRange
	}

	return bingPicture[start:end], nil
}
