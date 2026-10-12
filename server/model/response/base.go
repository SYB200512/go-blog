package response

type Captcha struct {
	CaptchaID string `json:"captcha_id"`
	// 验证码图片路径
	Picpath string `json:"picpath"`
}
