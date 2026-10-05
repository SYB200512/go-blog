package database

import (
	"server/global"
	"server/model/appTypes"
)

// Image 图片表
type Image struct {
	global.MODEL
	Name     string            `json:"name"`                       // 图片名称
	URL      string            `json:"url" gorm:"size:255;unique"` // 图片 URL
	Category appTypes.Category `json:"category"`                   // 图片分类
	Storage  appTypes.Storage  `json:"storage"`                    // 图片存储方式
}
