package service

import (
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"
)

type WebsiteService struct {
}

// WebsiteCarousel 获取轮播图
func (websiteService *WebsiteService) WebsiteCarousel() []string {
	var urls []string
	global.DB.Model(&database.Image{}).Where("category = ?", appTypes.Carousel).Pluck("url", &urls)
	return urls
}

// WebsiteNews 获取新闻
func (websiteService *WebsiteService) WebsiteNews(sourceStr string) (other.HotSearchData, error) {
	hotSearchData, err := ServiceGroupApp.HotSearchService.GetHotSearchDataBySource(sourceStr)
	if err != nil {
		return other.HotSearchData{}, err
	}
	return hotSearchData, nil
}

// WebsiteCalendar 获取日历
func (websiteService *WebsiteService) WebsiteCalendar(dateStr string) (other.Calendar, error) {
	calendar, err := ServiceGroupApp.CalendarService.GetCalendarByDate(dateStr)
	if err != nil {
		return other.Calendar{}, err
	}
	return calendar, nil
}

// WebsiteFooterLink 获取页底部链接
func (websiteService *WebsiteService) WebsiteFooterLink() []database.FooterLink {
	var footerLinks []database.FooterLink
	global.DB.Find(&footerLinks)
	return footerLinks
}

// WebsiteAddCarousel 添加轮播图
func (websiteService *WebsiteService) WebsiteAddCarousel(req request.WebsiteCarouselOperation) error {
	return utils.ChangeImagesCategory(global.DB, []string{req.Url}, appTypes.Carousel)
}

// WebsiteCancelCarousel 移除轮播图
func (websiteService *WebsiteService) WebsiteCancelCarousel(req request.WebsiteCarouselOperation) error {
	return utils.InitImagesCategory(global.DB, []string{req.Url})
}

// WebsiteCreateFooterLink 创建页脚链接
func (websiteService *WebsiteService) WebsiteCreateFooterLink(req database.FooterLink) error {
	return global.DB.Save(&req).Error
}

// WebsiteDeleteFooterLink 删除页脚链接
func (websiteService *WebsiteService) WebsiteDeleteFooterLink(req database.FooterLink) error {
	return global.DB.Delete(&req).Error
}
