package other

import (
	"server/model/request"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"gorm.io/gorm"
)

// MySQLOption MySQL 分页选项
type MySQLOption struct {
	request.PageInfo
	Order   string
	Where   *gorm.DB
	Preload []string
}

// EsOption Elasticsearch 分页选项
type EsOption struct {
	request.PageInfo
	Index          string
	Request        *search.Request
	SourceIncludes []string
}
