package main

type KfPerson struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	IdNumber string `db:"idcard"`
}

//缓存结果
type QueryResult struct {
	//开房者数据切片
	value []KfPerson
	//加入缓存的时间
	cacheTime int64
	//被查询的次数
	count int
}

func (qr *QueryResult) GetCacheTime() int64 {
	return qr.cacheTime
}