package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptHash 对密码进行 bcrypt 加密
// @param password 密码字符串
// @return 加密后的密码字符串
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
