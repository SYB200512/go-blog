package global

import (
	"server/config"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config      *config.Config
	Log         *zap.Logger
	DB          *gorm.DB
	EsClient    *elasticsearch.TypedClient
	RedisClient *redis.Client
	BlackCache  local_cache.Cache
)
