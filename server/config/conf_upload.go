package config

type Upload struct {
	Size int    `json:"size" yaml:"size"` //上传文件大小限制，单位：字节
	Path string `json:"path" yaml:"path"` //上传文件路径
}
