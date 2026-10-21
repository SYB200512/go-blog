package service

import (
	"context"
	"errors"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/elasticsearch"
	"server/model/other"
	"server/model/request"
	"server/utils"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"gorm.io/gorm"
)

type ArticleService struct {
}

// ArticleLike 点赞文章
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

// ArticleLikesList 获取用户点赞的文章列表
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
			{Match: map[string]types.MatchQuery{"tags": {Query: info.Tag}}},
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

func (articleService *ArticleService) ArticleCreate(req request.ArticleCreate) error {
	//1:判断文章标题是否存在
	b, err := articleService.Exits(req.Title)
	if err != nil {
		return err
	}
	if b {
		return errors.New("article title already exists")
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	articleToCreate := elasticsearch.Article{
		Cover:     req.Cover,
		Title:     req.Title,
		Category:  req.Category,
		Tags:      req.Tags,
		Abstract:  req.Abstract,
		Keyword:   req.Title,
		Content:   req.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	//创建了一篇文章同时要更新其他多个表的内容了
	//开启事务，要同时更新多个表，一个没有更新成功，其他表也不更新成功
	return global.DB.Transaction(func(tx *gorm.DB) error {
		//更新文章类别数量
		if err := articleService.UpdateCategoryCount(tx, "", articleToCreate.Category); err != nil {
			return err
		}
		//更新文章标签数量
		if err := articleService.UpdateTagsCount(tx, []string{}, articleToCreate.Tags); err != nil {
			return err
		}
		//更新图片表中的,添加图片时时为使用的，创建文章后要将图片类型改为封面
		if err := utils.ChangeImagesCategory(tx, []string{articleToCreate.Cover}, appTypes.Cover); err != nil {
			return err
		}

		//从文章内容中提取所有图片链接
		illustrations, err := utils.FindIllustrations(articleToCreate.Content)
		if err != nil {
			return err
		}

		//将文章中用的图片类别改为appTypes.Illustration
		if err := utils.ChangeImagesCategory(tx, illustrations, appTypes.Illustration); err != nil {
			return err
		}
		return articleService.Create(&articleToCreate)

	})
}
func (articleService *ArticleService) ArticleDelete(req request.ArticleDelete) error {
	//判断文章是否存在
	if len(req.IDs) == 0 {
		return nil
	}

	//开启事务，要同时删除多个表，一个没有删除成功，其他表也不删除成功
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			articleToDelete, err := articleService.Get(id)
			if err != nil {
				return err
			}
			if err := articleService.UpdateCategoryCount(tx, articleToDelete.Category, ""); err != nil {
				return err
			}
			if err := articleService.UpdateTagsCount(tx, articleToDelete.Tags, []string{}); err != nil {
				return err
			}
			if err := utils.InitImagesCategory(tx, []string{articleToDelete.Cover}); err != nil {
				return err
			}

			//将文章中用到的图片改为未使用
			illustrations, err := utils.FindIllustrations(articleToDelete.Content)
			if err != nil {
				return err
			}
			if err := utils.InitImagesCategory(tx, illustrations); err != nil {
				return err
			}

			//TODO同时删除该文章下的所有评论

		}

		return articleService.Delete(req.IDs)
	})
}
func (articleService *ArticleService) ArticleUpdate(req request.ArticleUpdate) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	articleToUpdate := struct {
		UpdateAt string   `json:"update_at"`
		Cover    string   `json:"cover"`
		Title    string   `json:"title"`
		Category string   `json:"category"`
		Tags     []string `json:"tags"`
		Abstract string   `json:"abstract"`
		Content  string   `json:"content"`
		Keyword  string   `json:"keyword"`
	}{
		UpdateAt: now,
		Cover:    req.Cover,
		Title:    req.Title,
		Category: req.Category,
		Tags:     req.Tags,
		Abstract: req.Abstract,
		Keyword:  req.Title,
		Content:  req.Content,
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {

		oldArticle, err := articleService.Get(req.ID)
		if err != nil {
			return err
		}
		//更新文章类别数量
		if err := articleService.UpdateCategoryCount(tx, oldArticle.Category, articleToUpdate.Category); err != nil {
			return err
		}
		//更新文章标签数量
		if err := articleService.UpdateTagsCount(tx, oldArticle.Tags, articleToUpdate.Tags); err != nil {
			return err
		}

		//更新图片表中的,添加图片时时为使用的，更新文章后要将图片类型改为封面
		if articleToUpdate.Cover != oldArticle.Cover {
			if err := utils.InitImagesCategory(tx, []string{oldArticle.Cover}); err != nil {
				return err
			}
			if err := utils.ChangeImagesCategory(tx, []string{articleToUpdate.Cover}, appTypes.Cover); err != nil {
				return err
			}
		}
		//从文章内容中提取所有图片链接
		oldIllustrations, err := utils.FindIllustrations(oldArticle.Content)
		if err != nil {
			return nil
		}
		//从新文章内容中提取所有图片链接
		newIllustrations, err := utils.FindIllustrations(articleToUpdate.Content)
		if err != nil {
			return nil
		}
		//对比新增和删除的图片
		addedILlustrations, removedIllustrations := utils.DiffArrays(oldIllustrations, newIllustrations)
		//将删除的图片类别改为未使用
		if err := utils.InitImagesCategory(tx, removedIllustrations); err != nil {
			return err
		}
		//将新增的图片类别改为appTypes.Illustration
		if err := utils.ChangeImagesCategory(tx, addedILlustrations, appTypes.Illustration); err != nil {
			return err
		}
		return articleService.Update(req.ID, articleToUpdate)
	})
}

// ArticleList 用于查询文章列表
func (articleService *ArticleService) ArticleList(info request.ArticleList) (list interface{}, total int64, err error) {
	req := &search.Request{
		Query: &types.Query{},
	}

	boolQuery := &types.BoolQuery{}

	// 根据标题查询
	if info.Title != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{Match: map[string]types.MatchQuery{"title": {Query: *info.Title}}})
	}

	// 根据简介查询
	if info.Abstract != nil {
		boolQuery.Must = append(boolQuery.Must, types.Query{Match: map[string]types.MatchQuery{"abstract": {Query: *info.Abstract}}})
	}

	// 根据类别筛选
	if info.Category != nil {
		boolQuery.Filter = []types.Query{
			{
				Term: map[string]types.TermQuery{
					"category": {Value: info.Category},
				},
			},
		}
	}

	// 根据条件执行查询
	if boolQuery.Must != nil || boolQuery.Filter != nil {
		req.Query.Bool = boolQuery
	} else {
		req.Query.MatchAll = &types.MatchAllQuery{}
		req.Sort = []types.SortCombinations{
			types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					"created_at": {Order: &sortorder.Desc},
				},
			},
		}
	}

	option := other.EsOption{
		PageInfo: info.PageInfo,
		Index:    elasticsearch.ArticleIndex(),
		Request:  req,
	}
	return utils.EsPagination(context.TODO(), option)
}
