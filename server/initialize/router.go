package initialize

import (
	"server/global"
	"server/router"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 设置gin模式
	gin.SetMode(global.Config.System.Env)
	Router := gin.Default()

	// 初始化路由组
	routerGroup := router.RouterGroupApp
	// 初始化基础路由
	publicGroup := Router.Group(global.Config.System.RouterPrefix)
	{
		//
		routerGroup.InitBaseRouter(publicGroup)
	}
	return Router
}
