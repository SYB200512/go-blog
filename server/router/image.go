package router

import (
	"server/api"

	"github.com/gin-gonic/gin"
)

// ImageRouter 图片路由
type ImageRouter struct {
}

// InitImageRouter 初始化图片路由
func (i *ImageRouter) InitImageRouter(Router *gin.RouterGroup) {
	imageRouter := Router.Group("image")

	imageApi := api.ApiGroupApp.ImageApi
	{
		imageRouter.POST("upload", imageApi.ImageUpload)
		imageRouter.DELETE("delete", imageApi.DeleteImage)
		imageRouter.GET("list", imageApi.GetImageList)
	}
}
