package service

import (
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
)

type WebsiteService struct {
}

// WebsiteCarousel 获取网站轮播图
func (websiteService *WebsiteService) WebsiteCarousel() []string {
	var urls []string
	//pluck：只取单个字段
	global.DB.Model(&database.Image{}).Where("category = ?", appTypes.Carousel).Pluck("url", (&urls))
	return urls

}

// WebsiteNews 获取网站新闻
func (websiteService *WebsiteService) WebsiteNews(sourceStr string) (other.HotSearchData, error) {
	//单独文件:hot_search.go
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

// WebsiteFooterLink 获取网站底部链接
func (websiteService *WebsiteService) WebsiteFooterLinkfunc() {

}

func (websiteService *WebsiteService) WebsiteAddCarousel() {

}
func (websiteService *WebsiteService) WebsiteCancelCarousel() {

}

// WebsiteCreateFooterLink 创建底部底部链接
func (websiteService *WebsiteService) WebsiteCreateFooterLink() {

}

// WebsiteDeleteFooterLink 删除底部底部链接
func (websiteService *WebsiteService) WebsiteDeleteFooterLink() {

}
