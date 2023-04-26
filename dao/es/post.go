package es

import (
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"mio-search/model"
	"reflect"
	"strconv"
)

const (
	PostIndex         = "post"
	PostTitle         = "title"
	PostContent       = "content"
	PostDeletedAt     = "deleted_at"
	PostEmpty         = ""
	PostHighlightPre  = "<em style=\"color:red\">"
	PostHighlightPost = "</em>"
)

func InsertPost(postList []*model.PostESEntity) (err error) {
	bulkRequest := es.Bulk().Index(PostIndex)
	for _, p := range postList {
		req := elastic.NewBulkIndexRequest().Doc(p)
		req.Id(strconv.FormatInt(p.PostId, 10))
		bulkRequest.Add(req)
	}
	_, err = bulkRequest.Do(ctx)
	return
}

func QueryPostByKeyword(p *model.SearchDTO) ([]*model.PostESEntity, error) {
	query := elastic.NewBoolQuery()
	query.MustNot(elastic.NewMatchQuery(PostTitle, PostEmpty))
	query.Should(elastic.NewMatchQuery(PostContent, p.Text), elastic.NewMatchQuery(PostTitle, p.Text))
	query.MustNot(elastic.NewExistsQuery(PostDeletedAt))

	highlight := elastic.NewHighlight()
	highlight.Fields(elastic.NewHighlighterField(PostContent), elastic.NewHighlighterField(PostTitle))
	highlight.PreTags(PostHighlightPre)
	highlight.PostTags(PostHighlightPost)

	result, err := es.Search().Index(PostIndex).From(p.Page - 1).Size(p.Size).Highlight(highlight).Query(query).Do(ctx)
	if err != nil {
		zap.L().Info("[dao es error] query post by keyword error ", zap.Error(err))
		return nil, err
	}

	var post model.PostESEntity
	var postIdList []*model.PostESEntity

	for index, element := range result.Each(reflect.TypeOf(post)) {
		var title, content string
		for _, str := range result.Hits.Hits[index].Highlight[PostTitle] {
			title += str
		}
		for _, str := range result.Hits.Hits[index].Highlight[PostContent] {
			content += str
		}
		entity := element.(model.PostESEntity)
		entity.Title = title
		entity.Content = content

		postIdList = append(postIdList, &entity)
	}

	return postIdList, nil
}
