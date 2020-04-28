package main

import "time"

// 声明一个CacheTimedData接口，该接口包含了一个获取缓存时间的方法
type CacheTimedData interface {
	GetCacheTime() int64 // 加入缓存的时间纳秒
}

//整理缓存
//删除加入最早的缓存
func UpdateCache(cacheMap *map[string]QueryResult) (delKey string){
	//预定义一个假设的最早时间
	earliestTime := time.Now().UnixNano()
	for key,value := range *cacheMap{
		if value.cacheTime < earliestTime{
			earliestTime = value.cacheTime
			delKey = key
		}
	}
	delete(*cacheMap,delKey)
	return delKey
}


// 缓存框架，传入一个实现了带缓存时间的数据接口
func UpdateCache2(cacheMap *map[string]CacheTimedData) (delKey string) {
	// 预设一个最早时间
	earliestTime := time.Now().UnixNano()
	for k, v := range *cacheMap {
		if v.GetCacheTime() < earliestTime {
			earliestTime = v.GetCacheTime()
			delKey = k
		}
	}
	delete(*cacheMap, delKey)
	return
}