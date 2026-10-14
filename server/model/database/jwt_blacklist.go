package database

import "server/global"

// JwtBlacklist jwt 黑名单表
type JwtBlacklist struct {
	global.MODEL
	Jwt string `json:"jwt" gorm:"type:text"` // JWT
}
