package flag

import (
	"server/global"
	"server/model/database"
)

// SQL 数据库表结构迁移，如果表不存在则创建；如果表存在则更新表结构
func SQL() error {
	return global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&database.Advertisement{},
		&database.ArticleCategory{},
		&database.ArticleLike{},
		&database.ArticleTag{},
		&database.Comment{},
		&database.Feedback{},
		&database.FooterLink{},
		&database.User{},
		&database.FriendLink{},
		&database.Image{},
		&database.JWTBlacklist{},
		&database.Login{},
	)
}
