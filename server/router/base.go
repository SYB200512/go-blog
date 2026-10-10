package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

// BaseRouter 基础路由
type BaseRouter struct {
}

// InitBaseRouter 初始化基础路由
func (b *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) {
	// 基础路由组
	baseRouter := Router.Group("base")

	baseApi := api.ApiGroupApp.BaseApi
	{
		baseRouter.POST("captcha", baseApi.Captcha)
		baseRouter.POST("sendEmailVerifivatonCode", baseApi.SendEmailVerificationCode)
		baseRouter.GET("qqLoginURL", baseApi.QQLoginURL)
	}
}
