package spider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mio-search/dao/mysql"
	"mio-search/model"
	"mio-search/util"
	"net/http"
)

func FetchInitPostList() error {
	url := "https://www.code-nav.cn/api/post/search/page/vo"
	var postList []*model.Post
	client := &http.Client{}

	for i := 1; i <= 5; i++ {
		payload := fmt.Sprintf("{\"current\":%d,\"pageSize\":8,\"sortField\":\"createTime\",\"sortOrder\":\"descend\",\"category\":\"文章\",\"reviewStatus\":1}", i)

		buf := bytes.NewBuffer([]byte(payload))

		newRequest, err := http.NewRequest(http.MethodPost, url, buf)
		newRequest.Header.Set("cookie", "输入你自己的cookie")
		newRequest.Header.Set("content-type", "application/json")

		resp, err := client.Do(newRequest)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		defer resp.Body.Close()

		// 获取响应的Body内容
		readAll, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("[io.ReadAll(response.Body) err]", err.Error())
			return err
		}

		data := make(map[string]interface{})

		err = json.Unmarshal(readAll, &data)
		if err != nil {
			fmt.Println("Unmarshal error ", err)
			return err
		}

		records := data["data"].(map[string]interface{})["records"].([]interface{})
		for _, record := range records {
			p := record.(map[string]interface{})

			post := &model.Post{
				Title:   p["title"].(string),
				Content: p["content"].(string),
				PostId:  util.GenSnowflakeID(),
				UserId:  133370088521728,
			}
			postList = append(postList, post)
		}
	}

	err := mysql.InsertPostList(postList)
	if err != nil {
		fmt.Println("insert post list error ", err)
		return err
	}

	return nil
}
