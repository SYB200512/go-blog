package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"
	"server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ArticleApi struct {
}

// ArticleLike 点赞文章
func (articleApi *ArticleApi) ArticleLike(c *gin.Context) {
	var req request.ArticleLike
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	err = articleService.ArticleLike(req)
	if err != nil {
		global.Log.Error("Failed to like article:", zap.Error(err))
		response.FailWithMessage("Failed to like article:", c)
		return
	}
	response.OkWithMessage("Article liked successfully", c)
}

// ArticleIsLike 判断文章是否被该用户点赞
func (articleApi *ArticleApi) ArticleIsLike(c *gin.Context) {
	var req request.ArticleLike
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	isLike, err := articleService.ArticleIsLike(req)
	if err != nil {
		global.Log.Error("Failed to get like status:", zap.Error(err))
		response.FailWithMessage("Failed to get like status:", c)
		return
	}
	response.OkWithData(isLike, c)
}

// ArticleLikesList 获取该用户点赞的文章列表
func (articleApi *ArticleApi) ArticleLikesList(c *gin.Context) {
	var pageInfo request.ArticleLikesList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	pageInfo.UserID = utils.GetUserID(c)
	list, total, err := articleService.ArticleLikesList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get article likes list:", zap.Error(err))
		response.FailWithMessage("Failed to get article likes list:", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

// ArticleInfoByID 获取文章详情
func (articleApi *ArticleApi) ArticleInfoByID(c *gin.Context) {
	var req request.ArticleInfoByID
	err := c.ShouldBindUri(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)

		return
	}
	article, err := articleService.ArticleInfoByID(req.ID)
	if err != nil {
		global.Log.Error("Failed to get article information:", zap.Error(err))
		response.FailWithMessage("Failed to get article information:", c)
		return
	}
	response.OkWithData(article, c)
}

// ArticleSearch 搜索文章
func (articleApi *ArticleApi) ArticleSearch(c *gin.Context) {
	var info request.ArticleSearch
	err := c.ShouldBindQuery(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := articleService.ArticleSearch(info)
	if err != nil {
		global.Log.Error("Failed to get article search results:", zap.Error(err))
		response.FailWithMessage("Failed to get article search results:", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

// ArticleCategory 获取所有文章类别及数量
func (articleApi *ArticleApi) ArticleCategory(c *gin.Context) {
	category, err := articleService.ArticleCategory()
	if err != nil {
		global.Log.Error("Failed to get article category:", zap.Error(err))
		response.FailWithMessage("Failed to get article category", c)
		return
	}
	response.OkWithData(category, c)
}

// ArticleTags 获取所有文章标签及数量
func (articleApi *ArticleApi) ArticleTags(c *gin.Context) {
	tags, err := articleService.ArticleTags()
	if err != nil {
		global.Log.Error("Failed to get article tags:", zap.Error(err))
		response.FailWithMessage("Failed to get article tags", c)
		return
	}
	response.OkWithData(tags, c)
}

func (articleApi *ArticleApi) ArticleCreate(c *gin.Context) {

}
func (articleApi *ArticleApi) ArticleDelete(c *gin.Context) {

}
func (articleApi *ArticleApi) ArticleUpdate(c *gin.Context) {

}

func (articleApi *ArticleApi) ArticleList(c *gin.Context) {

}
