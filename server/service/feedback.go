package service

import (
	"server/global"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"

	"github.com/gofrs/uuid"
)

type FeedbackService struct {
}

// FeedbackNew 获取最新反馈
func (feedbackService *FeedbackService) FeedbackNew() (feedbacks []database.Feedback, err error) {
	err = global.DB.Order("id desc").Limit(5).Find(&feedbacks).Error
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}

// FeedbackCreate 创建反馈
func (feedbackService *FeedbackService) FeedbackCreate(req request.FeedbackCreate) error {
	return global.DB.Create(&database.Feedback{
		UserUUID: req.UUID,
		Content:  req.Content,
	}).Error
}

// FeedbackInfo 获取用户反馈信息
func (feedbackService *FeedbackService) FeedbackInfo(uuid uuid.UUID) (feedbacks []database.Feedback, err error) {
	err = global.DB.Model(&database.Feedback{}).Order("id desc").Where("user_uuid = ?", uuid).Find(&feedbacks).Error
	if err != nil {
		return nil, err
	}
	return feedbacks, nil
}

// FeedbackDelete 删除反馈
func (feedbackService *FeedbackService) FeedbackDelete(req request.FeedbackDelete) error {
	if len(req.IDs) == 0 {
		return nil
	}
	return global.DB.Delete(&database.Feedback{}, req.IDs).Error
}

// FeedbackReply 回复反馈
func (feedbackService *FeedbackService) FeedbackReply(req request.FeedbackReply) error {
	return global.DB.Take(&database.Feedback{}, req.ID).Update("reply", req.Reply).Error
}

// FeedbackList 获取反馈列表
func (feedbackService *FeedbackService) FeedbackList(info request.PageInfo) (interface{}, int64, error) {
	option := other.MySQLOption{
		PageInfo: info,
	}

	return utils.MySQLPagination(&database.Feedback{}, option)
}
