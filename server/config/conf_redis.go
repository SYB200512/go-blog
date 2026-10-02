package config

type Redis struct {
	Address  string `json:"address" yaml:"address"`   //Redis地址，例如："127.0.0.1:6379"、"redis://127.0.0.1:6379"等
	Password string `json:"password" yaml:"password"` //Redis密码
	DB       int    `json:"db" yaml:"db"`             //Redis数据库索引，默认值为0
}
