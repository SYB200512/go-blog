package service

import (
	"server/utils"
	"time"

	"server/global"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type BaseService struct {
}

// SendEmailVerificationCode 发送邮箱验证码
func (baseService *BaseService) SendEmailVerificationCode(c *gin.Context, to string) error {
	// 生成验证码
	verificationCode := utils.GenerateVerificationCode(6)
	// 验证码过期时间
	expirationTime := time.Now().Add(5 * time.Minute).Unix()

	// 保存验证码到会话
	session := sessions.Default(c)
	session.Set("verification_code", verificationCode)
	session.Set("email", to)
	session.Set("expiration_time", expirationTime)
	_ = session.Save()

	// 发送验证码到邮箱
	subject := "您的邮箱验证码"
	body := `亲爱的用户[` + to + `]，<br/>
<br/>
感谢您注册` + global.Config.Website.Name + `的个人博客！为了确保您的邮箱安全，请使用以下验证码进行验证：<br/>
<br/>
验证码：[<font color="blue"><u>` + verificationCode + `</u></font>]<br/>
该验证码在 5 分钟内有效，请尽快使用。<br/>
<br/>
如果您没有请求此验证码，请忽略此邮件。
<br/>
如有任何疑问，请联系我们的支持团队：<br/>
邮箱：` + global.Config.Email.From + `<br/>
<br/>
祝好，<br/>` +
		global.Config.Website.Title + `<br/>
<br/>`

	_ = utils.Email(to, subject, body)

	return nil
}
