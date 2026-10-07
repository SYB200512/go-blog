package service

import (
	"context"
	"server/global"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type EsService struct {
}

// IndexCreate 创建索引
func (esService *EsService) IndexCreate(indexName string, mapping *types.TypeMapping) error {
	_, err := global.EsClient.Indices.Create(indexName).Mappings(mapping).Do(context.TODO())
	return err
}

// IndexDelete 删除索引
func (esService *EsService) IndexDelete(indexName string) error {
	_, err := global.EsClient.Indices.Delete(indexName).Do(context.TODO())
	return err
}

// IndexExists 检查索引是否存在
func (esService *EsService) IndexExists(indexName string) (bool, error) {
	return global.EsClient.Indices.Exists(indexName).Do(context.TODO())
}
