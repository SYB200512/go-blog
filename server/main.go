package main

import (
	"server/core"
	"server/flag"
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
	global.Redis = initialize.ConnectRedis()

	defer global.Redis.Close()
	// 初始化命令行参数
	flag.InitFlag()

	// 初始化定时任务
	initialize.InitCron()

	// 启动服务器
	core.RunServer()

}
