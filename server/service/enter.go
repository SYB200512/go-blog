package service

// ServiceGroup 服务组
type ServiceGroup struct {
	EsService
	BaseService
	JwtService
	GaodeService
	UserService
	QQService
	ImageService
	ArticleService
	CommentService
	FeedbackService
}

// ServiceGroupApp 服务组应用实例
var ServiceGroupApp = new(ServiceGroup)
