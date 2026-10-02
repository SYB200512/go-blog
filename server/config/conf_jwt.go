package config

type Jwt struct {
	AccessTokenSecret  string `json:"access_token_secret" yaml:"access_token_secret"`   //JWT密钥
	RefreshTokenSecret string `json:"refresh_token_secret" yaml:"refresh_token_secret"` //JWT密钥
	AccessTokenExpire  int64  `json:"access_token_expire" yaml:"access_token_expire"`   //JWT过期时间（秒）
	RefreshTokenExpire int64  `json:"refresh_token_expire" yaml:"refresh_token_expire"` //JWT过期时间（秒）
	Issuer             string `json:"issuer" yaml:"issuer"`                             //JWT发行者
}
