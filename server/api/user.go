package api

import (
	"errors"
	"server/global"
	"server/model/database"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserApi struct {
}

// Register 注册用户
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
	userApi.TokenNext(c, user)
}

// Login 登录用户
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

// EmailLogin 邮箱登录
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

// QQLogin QQ登录
func (userApi *UserApi) QQLogin(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage("Code is empty", c)
		return
	}

	accessTokenResponse, err := qqService.GetAccessTokenByCode(code)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user, err := userService.QQLogin(accessTokenResponse)
	if err != nil {
		global.Log.Error("Failed to login:", zap.Error(err))
		response.FailWithMessage("Failed to login", c)
		return
	}

	userApi.TokenNext(c, user)
}

// TokenNext 登录成功后，签发JWT token
func (userApi *UserApi) TokenNext(c *gin.Context, user database.User) {
	// 检查用户是否被冻结
	if user.Freeze {
		response.FailWithMessage("The user is frozen, contact the administrator", c)
		return
	}

	baseClaims := request.BaseClaims{
		UserID: user.ID,
		UUID:   user.UUID,
		RoleID: user.RoleID,
	}

	j := utils.NewJWT()

	// 创建访问令牌
	accessClaims := j.CreateAccessClaims(baseClaims)
	accessToken, err := j.CreateAccessToken(accessClaims)
	if err != nil {
		global.Log.Error("Failed to get accessToken:", zap.Error(err))
		response.FailWithMessage("Failed to get accessToken", c)
		return
	}

	// 创建刷新令牌
	refreshClaims := j.CreateRefreshClaims(baseClaims)
	refreshToken, err := j.CreateRefreshToken(refreshClaims)
	if err != nil {
		global.Log.Error("Failed to get refreshToken:", zap.Error(err))
		response.FailWithMessage("Failed to get refreshToken", c)
		return
	}

	// 是否开启了多地点登录拦截
	if !global.Config.System.UseMultipoint {
		// 设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
		return
	}

	// 检查 Redis 中是否已存在该用户的 JWT
	if jwtStr, err := jwtService.GetRedisJWT(user.UUID); errors.Is(err, redis.Nil) {
		// 不存在就设置新的
		if err := jwtService.SetRedisJWT(refreshToken, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage("Failed to set login status", c)
			return
		}

		// 设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
	} else if err != nil {
		// 出现错误处理
		global.Log.Error("Failed to set login status:", zap.Error(err))
		response.FailWithMessage("Failed to set login status", c)
	} else {
		// Redis 中已存在该用户的 JWT，将旧的 JWT 加入黑名单，并设置新的 token
		var blacklist database.JwtBlacklist
		blacklist.Jwt = jwtStr
		if err := jwtService.JoinInBlacklist(blacklist); err != nil {
			global.Log.Error("Failed to invalidate jwt:", zap.Error(err))
			response.FailWithMessage("Failed to invalidate jwt", c)
			return
		}

		// 设置新的 JWT 到 Redis
		if err := jwtService.SetRedisJWT(refreshToken, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage("Failed to set login status", c)
			return
		}

		// 设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
	}
}

// ForgotPassword 忘记密码,,找回密码
func (userApi *UserApi) ForgotPassword(c *gin.Context) {
	var req request.ForgotPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session := sessions.Default(c)

	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("Email not match", c)
		return
	}

	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Verification code not match", c)
		return
	}

	savedTime := session.Get("expiration_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("Verification code expired", c)
		return
	}

	err := userService.ForgotPassword(req)
	if err != nil {
		global.Log.Error("Failed to retrieve the password:", zap.Error(err))
		response.FailWithMessage("Failed to retrieve the password", c)
		return
	}
	response.OkWithMessage("Password retrieved successfully", c)

}

// 用户卡片信息
func (userApi *UserApi) UserCard(c *gin.Context) {
	var req request.UserCard
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userCard, err := userService.UserCard(req)
	if err != nil {
		global.Log.Error("Failed to retrieve user card:", zap.Error(err))
		response.FailWithMessage("Failed to retrieve user card", c)
		return
	}
	response.OkWithData(userCard, c)
}

// Logout 退出登录
func (userApi *UserApi) Logout(c *gin.Context) {
	userService.Logout(c)
	response.OkWithMessage("Logout successfully", c)
}

func (userApi *UserApi) UserResetPassword(c *gin.Context) {
	var req request.UserResetPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	err = userService.UserResetPassword(req)
	if err != nil {
		global.Log.Error("Failed to modify:", zap.Error(err))
		response.FailWithMessage("Failed tomodify,orginal password does not match the current password", c)
		return
	}
	response.OkWithMessage("Password modified successfully,please login again", c)
	userApi.Logout(c)
}

// 获取用户信息
func (userApi *UserApi) UserInfo(c *gin.Context) {
	userID := utils.GetUserID(c)
	user, err := userService.UserInfo(userID)
	if err != nil {
		global.Log.Error("Failed to get user information:", zap.Error(err))
		response.FailWithMessage("Failed to get user information", c)
		return
	}
	response.OkWithData(user, c)
}

// 修改用户信息
func (userApi *UserApi) UserChangeInfo(c *gin.Context) {
	var req request.UserChangeInfo
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	err = userService.UserChangeInfo(req)
	if err != nil {
		global.Log.Error("Failed to change user information:", zap.Error(err))
		response.FailWithMessage("Failed to change user information", c)
		return
	}
	response.OkWithMessage("User information modified successfully", c)
}
func (userApi *UserApi) UserWeather(c *gin.Context) {
	ip := c.ClientIP()
	weather, err := userService.UserWeather(c, ip)
	if err != nil {
		global.Log.Error("Failed to get weather:", zap.Error(err))
		response.FailWithMessage("Failed to get weather", c)
		return
	}
	response.OkWithData(weather, c)
}

func (userApi *UserApi) UserChart(c *gin.Context) {
	var req request.UserChart
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := userService.UserChart(req)
	if err != nil {
		global.Log.Error("Failed to get user chart:", zap.Error(err))
		response.FailWithMessage("Failed to get user chart", c)
		return
	}
	response.OkWithData(data, c)
}

// UserList 获取用户列表
func (userApi *UserApi) UserList(c *gin.Context) {
	var pageInfo request.UserList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := userService.UserList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get user list:", zap.Error(err))
		response.FailWithMessage("Failed to get user list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

// UserFreeze 冻结用户
func (userApi *UserApi) UserFreeze(c *gin.Context) {
	var req request.UserOperation
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.UserFreeze(req)
	if err != nil {
		global.Log.Error("Failed to freeze user:", zap.Error(err))
		response.FailWithMessage("Failed to freeze user", c)
		return
	}
	response.OkWithMessage("User frozen successfully", c)
}

func (userApi *UserApi) UserUnfreeze(c *gin.Context) {
	var req request.UserOperation
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.UserUnfreeze(req)
	if err != nil {
		global.Log.Error("Failed to unfreeze user:", zap.Error(err))
		response.FailWithMessage("Failed to unfreeze user", c)
		return
	}
	response.OkWithMessage("User unfrozen successfully", c)
}

// UserLoginList 获取用户登录记录
func (userApi *UserApi) UserLoginList(c *gin.Context) {
	var pageInfo request.UserLoginList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := userService.UserLoginList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get user login list:", zap.Error(err))
		response.FailWithMessage("Failed to get user login list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}
