package response

type Captcha struct {
	CaptchaID string `json:"captcha_id"`
	Picpath   string `json:"picpath"`
}
