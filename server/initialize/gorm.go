package initialize

import (
	"server/global"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm() *gorm.DB {
	//获取mysql配置
	mysqlCfg := global.Config.Mysql

	//打开mysql连接，初始化GORM
	db, err := gorm.Open(mysql.Open(mysqlCfg.Dsn()), &gorm.Config{
		Logger: logger.Default.LogMode(mysqlCfg.LogLevel()),
	})
	if err != nil {
		global.Log.Error("Failed to connect to mysql: %v", zap.Error(err))
	}
	//从gorm.DB获取原生mysql的sql.DB实例
	sqlDB, _ := db.DB()

	//设置mysql连接池的最大空闲连接数和最大打开连接数
	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConns)

	//返回gorm.DB实例
	return db

}
