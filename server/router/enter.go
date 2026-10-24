package router

// RouterGroup 路由组
type RouterGroup struct {
	BaseRouter
	UserRouter
	ImageRouter
	ArticleRouter
	CommentRouter
	FeedbackRouter
	AdvertisementRouter
	FriendLinkRouter
	WebsiteRouter
}

// RouterGroupApp 路由组应用实例
var RouterGroupApp = new(RouterGroup)
