package initialize

import (
	"os"
	"server/global"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

func ConnectES() *elasticsearch.TypedClient {
	// 从配置文件中获取 Elasticsearch 配置
	esCfg := global.Config.ES
	// 构建 Elasticsearch 配置
	cfg := elasticsearch.Config{
		Addresses: []string{esCfg.URL},
		Username:  esCfg.Username,
		Password:  esCfg.Password,
	}
	// 如果配置了在控制台打印日志，添加日志记录器
	if esCfg.IsConsolePrint {
		cfg.Logger = &elastictransport.ColorLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		}
	}

	// 创建 Elasticsearch 客户端
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		global.Log.Error("Failed to create Elasticsearch client", zap.Error(err))
		os.Exit(1)
	}
	return client
}
