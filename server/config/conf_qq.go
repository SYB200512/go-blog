package config

type QQ struct {
	Enable      bool   `json:"enable" yaml:"enable"`             //是否启用QQ登录
	AppID       string `json:"app_id" yaml:"app_id"`             //QQ应用ID
	AppSecret   string `json:"app_secret" yaml:"app_secret"`     //QQ应用密钥
	RedirectURI string `json:"redirect_uri" yaml:"redirect_uri"` //QQ应用回调URL

}

func (qq QQ) QQLoginURL() string {
	return "https://graph.qq.com/oauth2.0/authorize?" +
		"response_type=code&" +
		"client_id=" + qq.AppID + "&" +
		"redirect_uri=" + qq.RedirectURI

}
