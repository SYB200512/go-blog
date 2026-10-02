package config

type Es struct {
	URL            string `json:"url" yaml:"url"`                           //ES服务器地址
	Username       string `json:"username" yaml:"username"`                 //ES用户名或邮箱
	Password       string `json:"password" yaml:"password"`                 //ES密码或授权码
	IsConsolePrint bool   `json:"is_console_print" yaml:"is_console_print"` //是否在控制台打印ES日志
	// 其他ES配置...
}
