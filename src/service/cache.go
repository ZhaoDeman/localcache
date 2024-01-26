package service

import (
	"context"
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"log"
	"time"
)

type CallBack func() (value map[string]interface{}, err error)

// OnRemoveWithReason 当key被删除时进行回调
func OnRemoveWithReason(key string, entry []byte, reason bigcache.RemoveReason) {
	switch reason {
	case bigcache.Expired:
		log.Println(key, "超时")
		break
	case bigcache.NoSpace:
		log.Println("内存空间不足", key, "被删除")
		break
	case bigcache.Deleted:
		log.Println(key, "被删除")
		break
	}
	return
}

var TenantCache *ICache

func InitTenantCache(ctx context.Context, iCacheConfig ICacheConfig) {
	if TenantCache == nil {
		cache, err := bigcache.New(ctx, bigcache.Config{
			Shards:             256,
			LifeWindow:         24 * time.Hour,
			CleanWindow:        1 * time.Hour,
			MaxEntriesInWindow: 1000 * 10 * 60,
			MaxEntrySize:       500,
			StatsEnabled:       false,
			Verbose:            true,
			OnRemoveWithReason: OnRemoveWithReason,
			Hasher:             nil,
			HardMaxCacheSize:   0,
			Logger:             bigcache.DefaultLogger(),
		})
		if err != nil {
			log.Println("cache 初始化失败")
		}
		TenantCache = &ICache{
			Cache:        cache,
			ICacheConfig: iCacheConfig,
		}
	}
}

func (iCache *ICache) Get(key string) (value []byte, err error) {
	value, err = TenantCache.Cache.Get(key)
	if err != nil || len(value) == 0 {
		resMap := map[string]interface{}{}
		resMap, err = TenantCache.ICacheConfig.LoadFunc()
		value, err = json.Marshal(resMap[key])
		if len(value) > 0 {
			err = iCache.Set(key, value)
			if err != nil {
				log.Println("set 失败", err.Error())
			}
		}
		return
	}
	return
}

func (iCache *ICache) ReLoadFunc(ctx context.Context) {
	go func() {
		TenantCache.LoadFunc()
		for {
			select {
			case <-time.After(TenantCache.ICacheConfig.LoadInterval):
				TenantCache.LoadFunc()
			case <-ctx.Done():
				log.Println("ctx Done")
				break
			}
		}
	}()
}

func (iCache *ICache) LoadFunc() {
	resMap := map[string]interface{}{}
	var err error
	resMap, err = TenantCache.ICacheConfig.LoadFunc()
	if err != nil {
		log.Println("加载失败", err.Error())
	}
	for k, v := range resMap {
		entry, _ := json.Marshal(v)
		err = TenantCache.Set(k, entry)
		if err != nil {
			log.Println("set失败", err.Error())
		}
	}
}

func (iCache *ICache) Set(key string, entry []byte) error {
	return TenantCache.Cache.Set(key, entry)
}

