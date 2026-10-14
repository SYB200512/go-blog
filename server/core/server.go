package core

import (
	"server/global"
	"server/initialize"
	"server/service"

	"go.uber.org/zap"
	// "server/service"
)

type server interface {
	ListenAndServe() error
}

// RunServer 用于启动服务器
func RunServer() {
	// 初始化服务器地址
	addr := global.Config.System.Addr()
	// 初始化路由
	Router := initialize.InitRouter()

	// 加载所有的 JWT 黑名单，存入本地缓存
	service.LoadAll()

	// 初始化服务器并启动
	s := initServer(addr, Router)
	global.Log.Info("server run success on ", zap.String("address", addr))
	global.Log.Error(s.ListenAndServe().Error())
}
