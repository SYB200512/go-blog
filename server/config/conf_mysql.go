package config

import (
	"strconv"
	"strings"

	"gorm.io/gorm/logger"
)

type Mysql struct {
	Host         string `json:"host" yaml:"host"`                     //MySQL服务器地址
	Port         int    `json:"port" yaml:"port"`                     //MySQL服务器端口
	Config       string `json:"config" yaml:"config"`                 //MySQL连接字符串配置
	DBName       string `json:"db_name" yaml:"db_name"`               //MySQL数据库名称
	Username     string `json:"username" yaml:"username"`             //MySQL用户名或邮箱
	Password     string `json:"password" yaml:"password"`             //MySQL密码或授权码
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"` //MySQL最大空闲连接数
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"` //MySQL最大打开连接数
	LogMode      string `json:"log_mode" yaml:"log_mode"`             //MySQL日志级别
}

func (m Mysql) Dsn() string {
	//连接数据库的链接
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + strconv.Itoa(m.Port) + ")/" + m.DBName + "?" + m.LogMode
}

// 日志等级
// 就是将字段转化为实例
func (m Mysql) LogLevel() logger.LogLevel {
	switch strings.ToLower(m.LogMode) {
	case "info", "Info":
		return logger.Info

	case "error", "Error":
		return logger.Error

	case "silent", "Silent":
		return logger.Silent

	case "warn", "Warn":
		return logger.Warn

	default:
		return logger.Info
	}

}
