package database

import (
	"context"
	"server/global"
	"server/model/elasticsearch"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Comment 评论表
type Comment struct {
	global.MODEL
	ArticleID string    `json:"article_id"` // 文章 ID
	PID       *uint     `json:"p_id"`       // 父评论 ID
	PComment  *Comment  `json:"-" gorm:"foreignKey:PID"`
	Children  []Comment `json:"children" gorm:"foreignKey:PID"`                  // 子评论
	UserUUID  uuid.UUID `json:"user_uuid" gorm:"type:char(36)"`                  // 用户 uuid
	User      User      `json:"user" gorm:"foreignKey:UserUUID;references:UUID"` // 关联的用户
	Content   string    `json:"content"`                                         // 内容
}

// TODO 创建和删除评论时需要更新文章评论数
// 两个钩子函数，一个是创建评论时更新文章评论数，一个是删除评论时更新文章评论数

// 在gorm中创建了一条Comment记录后，gorm自动执行AfterCreate方法
func (c *Comment) AfterCreate(_ *gorm.DB) error {
	source := "ctx._source.comment += 1"

	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}

	_, err := global.EsClient.Update(elasticsearch.ArticleIndex(), c.ArticleID).Script(&script).Do(context.TODO())
	return err
}

// 在gorm中删除了一条Comment记录后，gorm自动执行BeforeDelete方法
func (c *Comment) BeforeDelete(_ *gorm.DB) error {

	var articleID string
	if err := global.DB.Model(&c).Pluck("article_id", &articleID).Error; err != nil {
		return err
	}

	source := "ctx._source.comments -= 1"

	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}

	_, err := global.EsClient.Update(elasticsearch.ArticleIndex(), articleID).Script(&script).Do(context.TODO())
	return err
}
