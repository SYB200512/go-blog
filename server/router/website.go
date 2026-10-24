package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type WebsiteRouter struct {
}

func (w *WebsiteRouter) InitWebsiteRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	websiteRouter := Router.Group("website")
	websitePublicRouter := PublicRouter.Group("website")

	websiteApi := api.ApiGroupApp.WebsiteApi
	{
		websiteRouter.POST("addCarousel", websiteApi.WebsiteAddCarousel)             //添加首页的轮播图
		websiteRouter.PUT("cancelCarousel", websiteApi.WebsiteCancelCarousel)        //取消首页的轮播图
		websiteRouter.POST("createFooterLink", websiteApi.WebsiteCreateFooterLink)   //创建页脚链接
		websiteRouter.DELETE("deleteFooterLink", websiteApi.WebsiteDeleteFooterLink) //删除页脚链接
	}
	{
		websitePublicRouter.GET("logo", websiteApi.WebsiteLogo)             //搜索条上显示的logo图片
		websitePublicRouter.GET("title", websiteApi.WebsiteTitle)           //搜索条上显示的标题
		websitePublicRouter.GET("info", websiteApi.WebsiteInfo)             //网站的配置信息，config.yaml中的website
		websitePublicRouter.GET("carousel", websiteApi.WebsiteCarousel)     //首页的轮播图，可以添加图片,前端可以设置默认图片，也可以自定义图片
		websitePublicRouter.GET("news", websiteApi.WebsiteNews)             //首页的新闻列表，显示最新的新闻一小时更新一次
		websitePublicRouter.GET("calendar", websiteApi.WebsiteCalendar)     //首页的日历，显示最新的事件一小时更新一次
		websitePublicRouter.GET("footerLink", websiteApi.WebsiteFooterLink) //页脚链接
	}
}
