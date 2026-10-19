package service

import (
	"server/model/elasticsearch"
	"server/model/request"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
)

type ArticleService struct {
}

func (articleService *ArticleService) ArticleLike() {

}
func (articleService *ArticleService) ArticleIsLike() {

}
func (articleService *ArticleService) ArticleLikesList() {

}

// ArticleInfoByID 用于通过ID从 Elasticsearch 获取文章
func (articleService *ArticleService) ArticleInfoByID(id string) (elasticsearch.Article, error) {
	//异步来增加浏览量
	go func() {
		articleView := articleService.NewArticleView()
		_ = articleView.Set(id)
	}()

	return articleService.Get(id)
}
func (articleService *ArticleService) ArticleSearch(info request.ArticleSearch) (interface{}, int64, error) {
	req := &search.Request{
		Query: &types.Query{},
	}

	boolQuery := &types.BoolQuery{}
	//设立boolQuery的原因就是info.Query和info.Tag、info.Category不为空
	if info.Query != "" {
		boolQuery.Should = []types.Query{ //带字段的匹配
			//模糊匹配标题
			{Match: map[string]types.MatchQuery{"title": {Query: info.Query}}},
			//精准匹配标题
			{Match: map[string]types.MatchQuery{"keyword": {Query: info.Query}}},

			{Match: map[string]types.MatchQuery{"abstract": {Query: info.Query}}},
			{Match: map[string]types.MatchQuery{"content": {Query: info.Query}}},
		}
	}

	if info.Tag != "" {
		boolQuery.Must = []types.Query{
			{Match: map[string]types.MatchQuery{"tag": {Query: info.Tag}}},
		}
	}

	if info.Category != "" {
		boolQuery.Filter = []types.Query{
			{Term: map[string]types.TermQuery{"category": {Value: info.Category}}},
		}
	}
	if boolQuery.Should != nil || boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
	}

	//设置排序字段
	if info.Sort != "" {
		var sortField string
		switch info.Sort {
		case "time":
			sortField = "create_at"
		case "view":
			sortField = "views"
		case "comment":
			sortField = "comments"
		case "like":
			sortField = "likes"
		default:
			sortField = "create_at"
		}

		var order sortorder.SortOrder
		if info.Order != "asc" {
			order = sortorder.Desc
		} else {
			order = sortorder.Asc
		}

		req.Sort = []types.SortCombination{
			types.SortOptions{
				SortOptions: map[string]types.Fieldsort{
					sortField: {Order: &order},
				},
			},
		}

	}
}
func (articleService *ArticleService) ArticleCategory() {

}
func (articleService *ArticleService) ArticleTags() {

}

func (articleService *ArticleService) ArticleCreate() {

}
func (articleService *ArticleService) ArticleDelete() {

}
func (articleService *ArticleService) ArticleUpdate() {

}
func (articleService *ArticleService) ArticleList() {

}
