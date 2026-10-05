package database

import "server/global"

// JWTBlacklist JWT 黑名单表
type JWTBlacklist struct {
	global.MODEL
	Jwt string `json:"jwt" gorm:"type:text"` // JWT
}
