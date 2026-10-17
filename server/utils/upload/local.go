package upload

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"server/global"
	"server/utils"
	"strings"
	"time"
)

type Local struct {
}

//实现OSS接口

func (*Local) UploadImage(file *multipart.FileHeader) (string, string, error) {
	size := float64(file.Size) / float64(1024*1024)
	//判断图片大小是否超过配置中的大小
	if size >= float64(global.Config.Upload.Size) {
		return "", "", fmt.Errorf("the image size exceeds the set size, the current size is: %.2f MB, the set size is: %d MB", size, global.Config.Upload.Size)

	}

	//先取出文件后缀，判断文件类型是否在白名单中
	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	if _, exists := WhiteImageList[ext]; !exists {
		return "", "", errors.New("don't upload files that aren't image types")
	}
	//拼接图片路径和文件名
	filename := utils.MD5V([]byte(name)) + "-" + time.Now().Format("20060102150405") + ext
	path := global.Config.Upload.Path + "/image/"

	//如果文件夹不存在，创建文件夹
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", "", err
	}
	//拼接图片路径
	filepath := path + filename
	//创建文件，写入图片内容
	out, err := os.Create(filepath)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	if _, err = io.Copy(out, f); err != nil {
		return "", "", err
	}

	return "/" + filepath, filename, nil
}

// 实现删除图片接口
func (*Local) DeleteImage(key string) error {
	path := global.Config.Upload.Path + "/image/" + key
	return os.Remove(path)
}
