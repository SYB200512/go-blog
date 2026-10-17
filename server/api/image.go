package api

import (
	"server/global"
	"server/model/request"
	"server/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ImageApi struct {
}

func (i *ImageApi) ImageUpload(c *gin.Context) {
	//接收表单数据中image字段
	_, header, err := c.Request.FormFile("image")
	if err != nil {
		global.Log.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return

	}
	//调用服务层上传图片
	url, err := imageService.ImageUpload(header)
	if err != nil {
		global.Log.Error("Failed to upload image:", zap.Error(err))
		response.FailWithMessage("Failed to upload image:", c)
		return
	}
	response.OkWithDetailed(response.ImageUpload{
		Url:     url,
		OssType: global.Config.System.OssType,
	}, "Image uploaded successfully:", c)
}

// DeleteImage 删除图片
func (i *ImageApi) DeleteImage(c *gin.Context) {
	var req request.ImageDelete
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = imageService.ImageDelete(req)
	if err != nil {
		global.Log.Error("Failed to delete image", zap.Error(err))
		response.FailWithMessage("Failed to delete image", c)
		return
	}
	response.OkWithMessage("Image deleted successfully", c)
}

func (i *ImageApi) GetImageList(c *gin.Context) {
	var pageInfo request.ImageList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := imageService.ImageList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get image list", zap.Error(err))
		response.FailWithMessage("Failed to get image list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)

}
