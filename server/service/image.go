package service

import (
	"mime/multipart"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"
	"server/utils/upload"

	"gorm.io/gorm"
)

type ImageService struct {
}

// ImageUpload 上传图片
func (imageApi *ImageService) ImageUpload(file *multipart.FileHeader) (string, error) {
	//上传图片到OSS
	oss := upload.NewOss()
	url, filename, err := oss.UploadImage(file)
	if err != nil {
		return "", err

	}
	//保存图片信息到数据库
	return url, global.DB.Create(&database.Image{
		Name:     filename,
		URL:      url,
		Category: appTypes.Null,
		Storage:  global.Config.System.Storage(),
	}).Error
}

func (imageApi *ImageService) ImageDelete(req request.ImageDelete) error {
	//先判断是否有图片ID
	if len(req.IDs) == 0 {
		return nil
	}

	var images []database.Image
	if err := global.DB.Find(&images, req.IDs).Error; err != nil {
		return err
	}
	for _, image := range images {
		if err := global.DB.Transaction(func(tx *gorm.DB) error {
			//创建删除文件工具：根据文件储存位置创建删除文件工具
			oss := upload.NewOssWithStorage(image.Storage)
			//删除数据库中的图片记录
			if err := global.DB.Delete(&image).Error; err != nil {
				return err
			}
			//删除真实文件
			return oss.DeleteImage(image.Name)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (imageService *ImageService) ImageList(info request.ImageList) (list interface{}, total int64, err error) {
	db := global.DB

	//根据名称查询
	if info.Name != nil {
		db = db.Where("name LIKE ?", "%"+*info.Name+"%")
	}

	//根据分类查询
	if info.Category != nil {
		category := appTypes.ToCategory(*info.Category)
		db = db.Where("category = ?", category)
	}

	//根据存储查询查询
	if info.Storage != nil {
		storage := appTypes.ToStorage(*info.Storage)
		db = db.Where("storage = ?", storage)
	}

	//根据分页查询
	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}
	return utils.MySQLPagination(&database.Image{}, option)
}
