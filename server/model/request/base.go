package request

// SendEmailVerificationCodeReq 发送邮箱验证码请求
type SendEmailVerificationCodeReq struct {
	Email     string `json:"email" binding:"required,email"`
	Captcha   string `json:"captcha" binding:"required,len=6"`
	CaptchaID string `json:"captcha_id" binding:"required"`
}
