package service

import (
	"errors"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type CommentService struct {
}

func (commentService *CommentService) CommentCreate(req request.CommentCreate) error {
	return global.DB.Create(&database.Comment{
		ArticleID: req.ArticleID,
		PID:       req.PID,
		UserUUID:  req.UserUUID,
		Content:   req.Content,
	}).Error
	//当创建一个评论时，需要更新文章的评论数量，应该是要用事务来处理，但是为了方便在database中comment.go中添加钩子函数

}

// CommentDelete 删除评论
func (commentService *CommentService) CommentDelete(c *gin.Context, req request.CommentDelete) error {
	//递归删除评论和子评论：写在了comment_helper中
	if len(req.IDs) == 0 {
		return nil
	}
	return global.DB.Transaction(func(tx *gorm.DB) error {
		for _, id := range req.IDs {
			var comment database.Comment
			if err := global.DB.Take(&comment, id).Error; err != nil {
				return err
			}
			//做一下权限认证
			userUUID := utils.GetUUID(c)
			userRoleID := utils.GetRoleID(c)
			if userUUID != comment.UserUUID && userRoleID != appTypes.Admin {
				return errors.New("you do not have permission to delete the comment")
			}
			if err := commentService.DeleteCommentAndChildren(tx, id); err != nil {
				return err
			}
		}
		return nil
	})

}

// CommentInfo 获取用户评论
func (commentService *CommentService) CommentInfo(uuid uuid.UUID) ([]database.Comment, error) {
	var rawComments []database.Comment
	if err := global.DB.Order("id desc").Where("user_uuid = ?", uuid).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid,username,avator,adress,signature")
	}).Find(&rawComments).Error; err != nil {
		return nil, err
	}
	for i := range rawComments {
		if err := commentService.LoadChildren(&rawComments[i]); err != nil {
			return nil, err
		}
	}

	//评论去重，如果当前评论的子评论存在为你的评论，就去除子评论。
	// 评论去重，如果当前评论的子评论存在为你的评论，就去除子评论，或者说当前评论的父评论存在为你的评论，就去除当前评论
	var comments []database.Comment

	idMap := commentService.FindChildCommentsIDByRootCommentUserUUID(rawComments)
	for i := range rawComments {
		if _, exists := idMap[rawComments[i].ID]; !exists {
			comments = append(comments, rawComments[i])
		}
	}
	return comments, nil
}

// CommentInfoByArticleID 根据文章ID查询所有评论，同时加载每个评论的子评论
func (commentService *CommentService) CommentInfoByArticleID(req request.CommentInfoByArticleID) ([]database.Comment, error) {
	var comments []database.Comment
	if err := global.DB.Where("article_id = ? AND p_id IS NULL ", req.ArticleID).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid,username,avator,adress,signature")
	}).Find(&comments).Error; err != nil {
		return nil, err
	}

	for i := range comments {
		if err := commentService.LoadChildren(&comments[i]); err != nil {
			return nil, err
		}
	}
	return comments, nil
}

// CommentNew 获取最新评论
func (commentService *CommentService) CommentNew() ([]database.Comment, error) {
	var comments []database.Comment
	err := global.DB.Order("id desc").Limit(5).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("uuid, username, avatar, address, signature")
	}).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// CommentList 获取评论列表
func (commentService *CommentService) CommentList(info request.CommentList) (interface{}, int64, error) {
	db := global.DB

	if info.ArticleID != nil {
		db = db.Where("article_id = ?", *info.ArticleID)
	}

	if info.UserUUID != nil {
		db = db.Where("user_uuid = ?", *info.UserUUID)
	}

	if info.Content != nil {
		db = db.Where("content LIKE ?", "%"+*info.Content+"%")
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.Comment{}, option)
}
