package service

import (
	"errors"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
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
func (userSercice *UserService) ForgotPassword() {

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
