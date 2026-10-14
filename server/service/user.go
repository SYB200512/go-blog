package service

import (
	"errors"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/utils"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type UserService struct {
}

// Register 注册用户
func (userSercice *UserService) Register(u database.User) (database.User, error) {
	//邮箱是否已注册
	if !errors.Is(global.DB.Where("email = ?", u.Email).First(&database.User{}).Error, gorm.ErrRecordNotFound) {
		return database.User{}, errors.New("this email address is already registered,please check the information you entered")
	}
	//密码加密
	u.Password = utils.BcryptHash(u.Password)
	//生成uuid
	u.UUID = uuid.Must(uuid.NewV4())
	u.Avatar = "/image/acatar.jpg"
	u.RoleID = appTypes.User
	u.Register = appTypes.Email

	if err := global.DB.Create(&u).Error; err != nil {
		return database.User{}, err
	}
	return u, nil
}

// EmailLogin 邮箱登录
func (userSercice *UserService) EmailLogin(u database.User) (database.User, error) {
	var user database.User
	err := global.DB.Where("email = ? ", u.Email).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return database.User{}, errors.New("incorrect email or password ")
		}
		return user, nil
	}
	return database.User{}, err
}

// QQLogin QQ登录
func (userSercice *UserService) QQLogin(accessTokenResponse other.AccessTokenResponse) (database.User, error) {
	var user database.User
	err := global.DB.Where("openid = ?", accessTokenResponse.Openid).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return database.User{}, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		userInfosponse, err := ServiceGroupApp.QQService.GetUserInfoByAccessTokenAndOpenid(accessTokenResponse.AccessToken, accessTokenResponse.Openid)
		if err != nil {
			return database.User{}, err
		}
		user.UUID = uuid.Must(uuid.NewV4())
		user.Username = userInfosponse.Nickname
		user.Openid = accessTokenResponse.Openid
		user.Avatar = userInfosponse.FigureurlQQ2
		user.RoleID = appTypes.User
		user.Register = appTypes.QQ
		if err := global.DB.Create(&user).Error; err != nil {
			return database.User{}, err
		}
	}
	return user, err
}

// ForgotPassword 忘记密码
func (userSercice *UserService) ForgotPassword(req request.ForgotPassword) error {
	var user database.User
	if err := global.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return err
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

func (userSercice *UserService) Logout() {

}
func (userSercice *UserService) UserResetPassword() {

}
func (userSercice *UserService) UserInfo() {

}
func (userSercice *UserService) UserChangeInfo() {

}
func (userSercice *UserService) UserWeather() {

}
func (userSercice *UserService) UserChart() {

}

func (userSercice *UserService) UserCard() {

}

func (userSercice *UserService) UserList() {

}
func (userSercice *UserService) UserFreeze() {

}
func (userSercice *UserService) UserUnfreeze() {

}
func (userSercice *UserService) UserLoginList() {

}
