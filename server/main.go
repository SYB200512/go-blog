package main

import (
	"server/core"
	"server/global"
	"server/initialize"
)

func main() {
	// 初始化配置
	global.Config = core.InitConf()
	// 初始化日志
	global.Log = core.InitLogger()

	initialize.OtherInit()
	// 初始化数据库
	global.DB = initialize.InitGorm()
	// 初始化 Elasticsearch 客户端
	global.EsClient = initialize.ConnectES()
	// 初始化 Redis 客户端
	global.RedisClient = initialize.ConnectRedis()

	core.RunServer()

}
