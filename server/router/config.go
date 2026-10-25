package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

type ConfigRouter struct {
}

func (c *ConfigRouter) InitConfigRouter(Router *gin.RouterGroup) {
	configRouter := Router.Group("config")

	configApi := api.ApiGroupApp.ConfigApi
	{
		//获取网站配置
		configRouter.GET("website", configApi.GetWebsite)
		//更新网站配置
		configRouter.PUT("website", configApi.UpdateWebsite)
		//获取系统配置
		configRouter.GET("system", configApi.GetSystem)
		configRouter.PUT("system", configApi.UpdateSystem)
		configRouter.GET("email", configApi.GetEmail)
		configRouter.PUT("email", configApi.UpdateEmail)
		configRouter.GET("qq", configApi.GetQQ)
		configRouter.PUT("qq", configApi.UpdateQQ)
		configRouter.GET("qiniu", configApi.GetQiniu)
		configRouter.PUT("qiniu", configApi.UpdateQiniu)
		configRouter.GET("jwt", configApi.GetJwt)
		configRouter.PUT("jwt", configApi.UpdateJwt)
		configRouter.GET("gaode", configApi.GetGaode)
		configRouter.PUT("gaode", configApi.UpdateGaode)
	}
}
