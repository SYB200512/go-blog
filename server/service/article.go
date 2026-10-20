package service

import (
	"context"
	"errors"
	"server/global"
	"server/model/database"
	"server/model/elasticsearch"
	"server/model/other"
	"server/model/request"
	"server/utils"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"gorm.io/gorm"
)

type ArticleService struct {
}

func (articleService *ArticleService) ArticleLike(req request.ArticleLike) error {
	//使用事务来确保数据的一致性
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var al database.ArticleLike
		var num int

		if errors.Is(tx.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).First(&al).Error, gorm.ErrRecordNotFound) {
			if err := tx.Create(&database.ArticleLike{UserID: req.UserID, ArticleID: req.ArticleID}).Error; err != nil {
				return err
			}
			num = 1
		} else {
			if err := tx.Delete(&al).Error; err != nil {
				return err
			}
			num = -1
		}

		//拼出ES的painless脚本，用于增加或减少点赞数
		source := "ctx._source.likes +=" + strconv.Itoa(num)
		//把这个ES的painless脚本转换为ES的UpdateRequest
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
		//对目标文章执行脚本更新：把 ES 里这篇文章的`likes` 字段，在原基础上加`num`
		_, err := global.EsClient.Update(elasticsearch.ArticleIndex(), req.ArticleID).Script(&script).Do(context.TODO())
		return err
	})
}

// ArticleIsLike 判断文章是否被该用户点赞
func (articleService *ArticleService) ArticleIsLike(req request.ArticleLike) (bool, error) {
	return !errors.Is(global.DB.Where("user_id = ? AND article_id = ?", req.UserID, req.ArticleID).First(&database.ArticleLike{}).Error, gorm.ErrRecordNotFound), nil
}

func (articleService *ArticleService) ArticleLikesList(info request.ArticleLikesList) (interface{}, int64, error) {
	db := global.DB.Where("user_id = ? ", info.UserID)
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	l, total, err := utils.MySQLPagination(&database.ArticleLike{}, option)
	if err != nil {
		return nil, 0, err
	}
	var list []struct {
		Id_     string                `json:"_id"`
		Source_ elasticsearch.Article `json:"_source"`
	}
	for _, articleLike := range l {
		article, err := articleService.Get(articleLike.ArticleID)
		if err != nil {
			return nil, 0, err
		}
		article.UpdatedAt = ""
		article.Keyword = ""
		article.Content = ""
		list = append(list, struct {
			Id_     string                `json:"_id"`
			Source_ elasticsearch.Article `json:"_source"`
		}{
			Id_:     articleLike.ArticleID,
			Source_: article,
		})
	}
	return list, total, nil
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

// ArticleSearch 用于搜索文章
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

		req.Sort = []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					sortField: {Order: &order},
				},
			},
		}
	}
	option := other.EsOption{
		PageInfo:       info.PageInfo,
		Index:          elasticsearch.ArticleIndex(),
		Request:        req,
		SourceIncludes: []string{"created_at", "cover", "title", "abstract", "tags", "views", "comments", "likes"},
	}
	return utils.EsPagination(context.TODO(), option)

}

// ArticleCategory 获取所有文章类别及数量
func (articleService *ArticleService) ArticleCategory() ([]database.ArticleCategory, error) {
	var category []database.ArticleCategory
	if err := global.DB.Find(&category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// ArticleTags 获取所有文章标签及数量
func (articleService *ArticleService) ArticleTags() ([]database.ArticleTag, error) {
	var tags []database.ArticleTag
	if err := global.DB.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (articleService *ArticleService) ArticleCreate() {

}
func (articleService *ArticleService) ArticleDelete() {

}
func (articleService *ArticleService) ArticleUpdate() {

}
func (articleService *ArticleService) ArticleList() {

}
