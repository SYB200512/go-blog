package api

import (
	"net/http"
	"server/global"
	"server/model/response"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WebsiteApi struct {
}

func (websiteApi *WebsiteApi) WebsiteLogo(c *gin.Context) {
	//返回搜索条上显示的logo图片:url设置在config.yaml中的website.logo

	//设置重定向到logo图片的url
	if global.Config.Website.Logo != "" {
		c.Redirect(http.StatusMovedPermanently, global.Config.Website.Logo)
	} else {
		c.Redirect(http.StatusMovedPermanently, "/image/logo.png")
	}
}

// Title 获取网站标题
func (websiteApi *WebsiteApi) WebsiteTitle(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"title": global.Config.Website.Title})
}

// WebsiteInfo 获取网站信息
func (websiteApi *WebsiteApi) WebsiteInfo(c *gin.Context) {

	response.OkWithData(global.Config.Website, c)
}

// WebsiteCarousel 获取网站轮播图
func (websiteApi *WebsiteApi) WebsiteCarousel(c *gin.Context) {
	urls := websiteService.WebsiteCarousel()
	response.OkWithData(urls, c)
}

// WebsiteNews 获取网站新闻
func (websiteApi *WebsiteApi) WebsiteNews(c *gin.Context) {
	sourceStr := c.Query("source")
	hotSearchData, err := websiteService.WebsiteNews(sourceStr)
	if err != nil {
		global.Log.Error("Failed to get news:", zap.Error(err))
		response.FailWithMessage("Failed to get news", c)
		return

	}
	response.OkWithData(hotSearchData, c)
}

// WebsiteCalendar 获取日历
func (website *WebsiteApi) WebsiteCalendar(c *gin.Context) {
	dateStr := time.Now().Format("2006/0102")
	calendar, err := websiteService.WebsiteCalendar(dateStr)
	if err != nil {
		global.Log.Error("Failed to get calendar:", zap.Error(err))
		response.FailWithMessage("Failed to get calendar", c)
		return
	}
	response.OkWithData(calendar, c)
}

// WebsiteFooterLink 获取网站底部链接
func (websiteApi *WebsiteApi) WebsiteFooterLink(c *gin.Context) {

}

func (websiteApi *WebsiteApi) WebsiteAddCarousel(c *gin.Context) {

}
func (websiteApi *WebsiteApi) WebsiteCancelCarousel(c *gin.Context) {
}

func (websiteApi *WebsiteApi) WebsiteCreateFooterLink(c *gin.Context) {
}

func (websiteApi *WebsiteApi) WebsiteDeleteFooterLink(c *gin.Context) {
}
