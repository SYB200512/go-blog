package service

import (
	"context"
	"server/global"
	"strconv"
)

// NewArticleView 用于创建文章浏览量统计数据库
func (articleService *ArticleService) NewArticleView() CountDB {
	return CountDB{
		Index: "article_views",
	}
}

type CountDB struct {
	Index string
}

// Set 在原有基础上加一
func (c CountDB) Set(id string) error {
	num, _ := global.Redis.HGet(context.Background(), c.Index, id).Int()
	num++
	err := global.Redis.HSet(context.Background(), c.Index, id, num).Err()
	return err
}

// GetInfo 取出数据
func (c CountDB) GetInfo() map[string]int {
	var Info = map[string]int{}
	maps := global.Redis.HGetAll(context.Background(), c.Index).Val()
	for id, val := range maps {
		num, _ := strconv.Atoi(val)
		Info[id] = num
	}
	return Info
}

// Clear 清除数据
func (c CountDB) Clear() {
	global.Redis.Del(context.Background(), c.Index)
}
