package config

type Gaode struct {
	Enable bool   `json:"enable" yaml:"enable"`   //是否启用高德地图
	Key    string `json:"api_key" yaml:"api_key"` //高德地图API密钥

}
