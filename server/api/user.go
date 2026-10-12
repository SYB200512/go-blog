package api

import (
	"server/model/request"

	"github.com/gin-gonic/gin"
)

type UserApi struct {
}

func (userApi *UserApi) Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBindJSON(&req)
	if err != nil {

		return
	}
}
func (userApi *UserApi) Login(c *gin.Context) {

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
