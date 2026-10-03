package main

import (
	"server/core"
	"server/global"
)

func main() {
	global.Config = core.InitConfig()
	global.Log = core.InitLogger()

	core.RunServer()

}
