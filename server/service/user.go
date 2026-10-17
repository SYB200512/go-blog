package service

import (
	"errors"
	"fmt"
	"server/global"
	"server/model/appTypes"
	"server/model/database"
	"server/model/other"
	"server/model/request"
	"server/model/response"
	"server/utils"
	"time"

	"github.com/gin-gonic/gin"
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

// 用户卡片信息
func (userSercice *UserService) UserCard(req request.UserCard) (response.UserCard, error) {
	var user database.User
	if err := global.DB.Where("uuid = ?", req.UUID).Select("uuid", "username", "avator", "address", "signature").First(&user).Error; err != nil {
		return response.UserCard{}, err
	}
	return response.UserCard{
		UUID:      user.UUID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Address:   user.Address,
		Signature: user.Signature,
	}, nil
}

// Logout 退出登录
func (userSercice *UserService) Logout(c *gin.Context) {
	uuid := utils.GetUUID(c)
	jwtStr := utils.GetRefreshToken(c)
	utils.ClearRefreshToken(c)
	global.Redis.Del(c, uuid.String())
	_ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{
		Jwt: jwtStr,
	})

}

// UserResetPassword 重置密码
func (userSercice *UserService) UserResetPassword(req request.UserResetPassword) error {
	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	if ok := utils.BcryptCheck(req.Password, user.Password); !ok {
		return errors.New("orginal password does not match the current password")
	}
	user.Password = utils.BcryptHash(req.NewPassword)
	return global.DB.Save(&user).Error
}

// UserInfo 获取用户信息
func (userSercice *UserService) UserInfo(userID uint) (database.User, error) {
	var user database.User
	if err := global.DB.Take(&user, userID).Error; err != nil {
		return database.User{}, err
	}
	return user, nil
}

// UserChangeInfo 修改用户信息
func (userSercice *UserService) UserChangeInfo(req request.UserChangeInfo) error {

	var user database.User
	if err := global.DB.Take(&user, req.UserID).Error; err != nil {
		return err
	}
	return global.DB.Model(&user).Updates(req).Error
}

func (userSercice *UserService) UserWeather(c *gin.Context, ip string) (string, error) {
	//防止频繁刷新天气，消耗api次数，先从redis中获取天气信息
	result, err := global.Redis.Get(c, "weather-"+ip).Result()
	if err != nil {
		ipResponse, err := ServiceGroupApp.GaodeService.GetLocationByIP(ip)
		if err != nil {
			return "", err
		}
		live, err := ServiceGroupApp.GaodeService.GetWeatherByAdcode(ipResponse.Adcode)
		if err != nil {
			return "", err
		}
		weather := "地区：" + live.Province + "-" + live.City + " 天气：" + live.Weather + " 温度：" + live.Temperature + "°C" + " 风向：" + live.WindDirection + " 风级：" + live.WindPower + " 湿度：" + live.Humidity + "%"

		// 将天气数据存入redis
		if err := global.Redis.Set(c, "weather-"+ip, weather, time.Hour*1).Err(); err != nil {
			return "", err
		}
		return weather, nil
	}
	return result, nil
}

func (userSercice *UserService) UserChart(req request.UserChart) (response.UserChart, error) {
	Where := global.DB.Where(fmt.Sprintf("date_sub(curdate(),interval %d day) <= created_at", req.Date))

	var res response.UserChart

	startDate := time.Now().AddDate(0, 0, -req.Date)
	for i := 1; i < req.Date; i++ {
		res.DataList = append(res.DataList, startDate.AddDate(0, 0, i).Format("2006-01-02"))
	}

	loginCounts := utils.FetchDateCounts(global.DB.Model(&database.Login{}), Where)
	registerCounts := utils.FetchDateCounts(global.DB.Model(&database.User{}), Where)
	for _, date := range res.DataList {
		loginCount := loginCounts[date]
		registerCount := registerCounts[date]
		res.LoginData = append(res.LoginData, loginCount)
		res.RegisterData = append(res.RegisterData, registerCount)
	}
	return res, nil
}

// UserList 获取用户列表
func (userService *UserService) UserList(info request.UserList) (interface{}, int64, error) {
	db := global.DB

	if info.UUID != nil {
		db = db.Where("uuid = ?", info.UUID)
	}

	if info.Register != nil {
		db = db.Where("register = ?", info.Register)
	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
	}

	return utils.MySQLPagination(&database.User{}, option)

}

// UserFreeze 冻结用户
func (userSercice *UserService) UserFreeze(req request.UserOperation) error {
	var user database.User
	if err := global.DB.Take(&user, req.ID).Update("freeze", true).Error; err != nil {
		return err
	}
	jwtStr, _ := ServiceGroupApp.JwtService.GetRedisJWT(user.UUID)
	if jwtStr != "" {
		_ = ServiceGroupApp.JwtService.JoinInBlacklist(database.JwtBlacklist{
			Jwt: jwtStr,
		})
	}
	return nil
}
func (userSercice *UserService) UserUnfreeze(req request.UserOperation) error {
	return global.DB.Take(&database.User{}, req.ID).Update("freeze", false).Error
}

// 获取用户登录记录
func (userSercice *UserService) UserLoginList(info request.UserLoginList) (interface{}, int64, error) {
	db := global.DB

	if info.UUID != nil {
		var userID uint
		if err := global.DB.Model(&database.User{}).Where("uuid = ?", info.UUID).Pluck("id", &userID).Error; err != nil {
			return nil, 0, err
		}
		db.Where("user_id = ?", userID)

	}

	option := other.MySQLOption{
		PageInfo: info.PageInfo,
		Where:    db,
		Preload:  []string{"User"},
	}

	return utils.MySQLPagination(&database.Login{}, option)
}
