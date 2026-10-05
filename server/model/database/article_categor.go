package database

type ArticleCategory struct {
	Category string `json:"category" gorm:"primaryKey"` // 文章分类
	Number   int    `json:"number"`                     // 文章数量
}
