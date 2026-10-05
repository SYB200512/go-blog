package database

type ArticleTag struct {
	Tag    string `json:"tag" gorm:"primaryKey"` // 文章标签
	Number int    `json:"number"`                // 文章数量
}
