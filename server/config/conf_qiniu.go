package config

type Qiniu struct {
	Zone         string `json:"zone" yaml:"zone"`                     //七牛云区域，例如："z0"、"z1"、"z2"等
	Bucket       string `json:"bucket" yaml:"bucket"`                 //七牛云存储桶名称
	ImgPath      string `json:"img_path" yaml:"img_path"`             //七牛云图片路径
	AccessKey    string `json:"access_key" yaml:"access_key"`         //七牛云访问密钥
	SecretKey    string `json:"secret_key" yaml:"secret_key"`         //七牛云密钥
	UseHTTPS     bool   `json:"use_https" yaml:"use_https"`           //是否使用HTTPS协议
	UseCdnDomain bool   `json:"use_cdn_domain" yaml:"use_cdn_domain"` //是否使用CDN域名访问

}
