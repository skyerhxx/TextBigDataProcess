package main

import "time"

type TimedData interface {
	//获得加入缓存的时间纳秒
	GetCacheTime() int64
}


//整理缓存
//删除加入最早的缓存
func UpdateCache(cacheMap *map[string]TimedData) (delKey string){
	//预定义一个假设的最早时间
	earliestTime := time.Now().UnixNano()
	for key,value := range *cacheMap{
		if value.GetCacheTime() < earliestTime{
			earliestTime = value.GetCacheTime()
			delKey = key
		}
	}
	delete(*cacheMap,delKey)
	return delKey
}