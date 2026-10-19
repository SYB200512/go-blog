package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type ArticleRouter struct {
}

func (a *ArticleRouter) InitArticleRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup, AdminRouter *gin.RouterGroup) {
	articleRouter := Router.Group("article")
	articlePublicRouter := PublicRouter.Group("article")
	articleAdminRouter := AdminRouter.Group("article")

	articleApi := api.ApiGroupApp.ArticleApi
	{
		//文章收藏
		articleRouter.POST("like", articleApi.ArticleLike)
		//判断是否收藏
		articleRouter.GET("isLike", articleApi.ArticleIsLike)
		//收藏列表
		articleRouter.GET("likesList", articleApi.ArticleLikesList)
	}
	{
		//文章详情
		articlePublicRouter.GET(":id", articleApi.ArticleInfoByID)
		//文章搜索
		articlePublicRouter.GET("search", articleApi.ArticleSearch)
		//文章分类
		articlePublicRouter.GET("category", articleApi.ArticleCategory)
		//文章标签
		articlePublicRouter.GET("tags", articleApi.ArticleTags)
	}
	{
		//文章创建
		articleAdminRouter.POST("create", articleApi.ArticleCreate)
		//文章删除
		articleAdminRouter.DELETE("delete", articleApi.ArticleDelete)
		//文章更新
		articleAdminRouter.PUT("update", articleApi.ArticleUpdate)
		//文章列表更新
		articleAdminRouter.PUT("update", articleApi.ArticleUpdate)
		//文章列表
		articleAdminRouter.GET("list", articleApi.ArticleList)
	}
}
