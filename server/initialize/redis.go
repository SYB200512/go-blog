package initialize

import (
	"context"
	"server/global"

	"os"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func ConnectRedis() *redis.Client {
	//获取redis配置
	redisCfg := global.Config.Redis

	//创建redis.Client实例
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Address,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	//测试redis连接
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		//记录错误日志
		global.Log.Error("Failed to connect to redis", zap.Error(err))
		//退出程序
		os.Exit(1)
	}
	return client
}
