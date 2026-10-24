package service

import (
	"context"
	"encoding/json"
	"server/global"
	"server/model/other"
	"server/utils/hotSearch"
	"time"
)

type HotSearchService struct {
}

// GetHotSearchDataBySource 获取热门搜索数据
func (hotSearchService *HotSearchService) GetHotSearchDataBySource(sourceStr string) (other.HotSearchData, error) {
	result, err := global.Redis.Get(context.TODO(), sourceStr).Result()
	if err != nil {
		source := hotSearch.NewSource(sourceStr)
		hotSearchData, err := source.GetHotSearchData(30)
		if err != nil {
			return other.HotSearchData{}, err
		}
		bytes, err := json.Marshal(hotSearchData)
		if err != nil {
			return other.HotSearchData{}, err
		}
		if err := global.Redis.Set(context.TODO(), sourceStr, bytes, time.Hour).Err(); err != nil {
			return other.HotSearchData{}, err
		}
		return hotSearchData, nil
	}
	var hotSearchData other.HotSearchData
	if err := json.Unmarshal([]byte(result), &hotSearchData); err != nil {
		return other.HotSearchData{}, err
	}
	return hotSearchData, nil
}
