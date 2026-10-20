package request

type ArticleInfoByID struct {
	ID string `json:"id" form:"id" uri:"id" binding:"required"`
}

type ArticleSearch struct {
	Query    string `json:"query" form:"query" uri:"query" binding:"max=50"`
	Category string `json:"category" form:"category" uri:"category"`
	Tag      string `json:"tag" form:"tag" uri:"tag"`
	Sort     string `json:"sort" form:"sort" uri:"sort"`
	Order    string `json:"order" form:"order" uri:"order" binding:"required"` //asc or desc
	PageInfo
}

type ArticleLike struct {
	UserID    uint   `json:"_"`
	ArticleID string `json:"article_id" form:"article_id" binding:"required"`
}

type ArticleLikesList struct {
	UserID uint `json:"_"`
	PageInfo
}
