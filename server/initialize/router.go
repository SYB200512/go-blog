package initialize

import (
	"net/http"
	"server/global"
	"server/middleware"
	"server/router"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 设置gin模式
	gin.SetMode(global.Config.System.Env)
	Router := gin.Default()

	//使用日志记录中间件
	Router.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 使用gin会话路由
	var store = cookie.NewStore([]byte(global.Config.System.SessionsSecret))
	// 初始化会话中间件
	Router.Use(sessions.Sessions("session", store))

	// "uploads" 是URL路径前缀，http.Dir("uploads")是实际文件系统中存储文件的目录
	//把收到的 URL 请求（`/uploads/xxx` ）映射到磁盘文件（`uploads/xxx` ），让这个 URL 真正能访问到图片字节
	Router.StaticFS(global.Config.Upload.Path, http.Dir(global.Config.Upload.Path))

	// 创建路由组
	routerGroup := router.RouterGroupApp
	// 初始化公共路由
	publicGroup := Router.Group(global.Config.System.RouterPrefix)
	// 初始化私有路由
	privateGroup := Router.Group(global.Config.System.RouterPrefix)
	privateGroup.Use(middleware.JWTAuth())
	// 初始化管理员路由
	adminGroup := Router.Group(global.Config.System.RouterPrefix)
	adminGroup.Use(middleware.JWTAuth()).Use(middleware.AdminAuth())

	{
		routerGroup.InitBaseRouter(publicGroup)
	}
	{
		routerGroup.InitUserRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitArticleRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitCommentRouter(privateGroup, publicGroup, adminGroup)
		routerGroup.InitFeedbackRouter(privateGroup, publicGroup, adminGroup)

	}
	{
		routerGroup.InitImageRouter(adminGroup)
		routerGroup.InitAdvertisementRouter(adminGroup, publicGroup)
		routerGroup.InitFriendLinkRouter(adminGroup, publicGroup)
		routerGroup.InitWebsiteRouter(adminGroup, publicGroup)
	}

	return Router
}
