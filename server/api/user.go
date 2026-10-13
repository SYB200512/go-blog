package api

import (
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserApi struct {
}

func (userApi *UserApi) Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//邮箱验证
	session := sessions.Default(c)
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("This email doesn't match the email in the verification code", c)
		return
	}
	saveCode := session.Get("code")
	if saveCode == nil || saveCode.(string) != req.VerificationCode {
		response.FailWithMessage("This verification code is not correct", c)
		return
	}

	savedTime := session.Get("expire_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("This verification code is expired,please resent it", c)
		return
	}

	u := database.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	user, err := userService.Register(u)
	if err != nil {
		global.Log.Error("Failed to register user:", zap.Error(err))
		response.FailWithMessage("Failed to register user", c)
		return
	}

	//注册成功后，自动登录，签发JWT token
	//userApi.TokenNext(c, user)
}

// 登录
func (userApi *UserApi) Login(c *gin.Context) {
	switch c.Query("flag") {
	case "email":
		userApi.EmailLogin(c)
	case "qq":
		userApi.QQLogin(c)
	default:
		userApi.EmailLogin(c)
	}
}

// 邮箱登录
func (userApi *UserApi) EmailLogin(c *gin.Context) {
	var req request.Login
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//验证码验证
	if store.Verify(req.CaptchaID, req.Captcha, true) {
		u := database.User{
			Email:    req.Email,
			Password: req.Password,
		}
		user, err := userService.EmailLogin(u)
		if err != nil {
			global.Log.Error("Failed to login user:", zap.Error(err))
			response.FailWithMessage("Failed to login user", c)
			return
		}
		userApi.TokenNext(c, user)
		return
	}

	response.FailWithMessage("Captcha is not correct", c)
	return
}

func (userApi *UserApi) QQLogin(c *gin.Context) {

}

func (userApi *UserApi) Logout(c *gin.Context) {

}
func (userApi *UserApi) UserResetPassword(c *gin.Context) {

}
func (userApi *UserApi) UserInfo(c *gin.Context) {

}
func (userApi *UserApi) UserChangeInfo(c *gin.Context) {

}
func (userApi *UserApi) UserWeather(c *gin.Context) {

}
func (userApi *UserApi) UserChart(c *gin.Context) {

}
func (userApi *UserApi) ForgotPassword(c *gin.Context) {

}
func (userApi *UserApi) UserCard(c *gin.Context) {

}

func (userApi *UserApi) UserList(c *gin.Context) {

}
func (userApi *UserApi) UserFreeze(c *gin.Context) {

}
func (userApi *UserApi) UserUnfreeze(c *gin.Context) {

}
func (userApi *UserApi) UserLoginList(c *gin.Context) {

}
