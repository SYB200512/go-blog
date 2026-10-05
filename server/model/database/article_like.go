package database

import "server/global"

type ArticleLike struct {
	global.MODEL
	ArticleID string `json:"article_id"` // 文章 ID
	UserID    uint   `json:"user_id"`    // 用户 ID
	User      User   `json:"_" gorm:"foreignKey:UserID"`
}
