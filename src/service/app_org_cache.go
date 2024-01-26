package service


import (
	"context"
	"encoding/json"
	"github.com/allegro/bigcache/v3"
	"log"
	"time"
)

var AppCache *ICache
var OrganizationCache *ICache
var ChannelCache *ICache

type ICache struct {
	Cache        *bigcache.BigCache
	ICacheConfig ICacheConfig
	AllKey       map[string]struct{}
}

type ICacheConfig struct {
	LoadFunc     CallBack
	CacheName    string
	LoadInterval time.Duration
	bConfig      *bigcache.Config
}

func InitAppCache(ctx context.Context, iCacheConfig ICacheConfig) {
	if AppCache == nil {
		cache, err := GetDefaultBigCache(ctx)
		if err != nil {
			log.Println("应用cache 初始化失败")
		}
		AppCache = &ICache{
			Cache:        cache,
			ICacheConfig: iCacheConfig,
		}
	}
}

func GetDefaultBigCache(ctx context.Context) (cache *bigcache.BigCache, err error) {
	cache, err = bigcache.New(ctx, bigcache.Config{
		Shards:             1024,
		LifeWindow:         24 * time.Hour,
		CleanWindow:        1 * time.Hour,
		MaxEntriesInWindow: 1000 * 10 * 15,
		MaxEntrySize:       256,
		StatsEnabled:       false,
		Verbose:            true,
		OnRemoveWithReason: OnRemoveWithReason,
		Hasher:             nil,
		HardMaxCacheSize:   0,
		Logger:             bigcache.DefaultLogger(),
	})
	return
}

func GetNewBigCache(ctx context.Context, config *bigcache.Config) (cache *bigcache.BigCache, err error) {
	if config == nil {
		return GetDefaultBigCache(ctx)
	}
	cache, err = bigcache.New(ctx, *config)
	return
}

func InitOrganizationCache(ctx context.Context, iCacheConfig ICacheConfig) {
	if OrganizationCache == nil {
		cache, err := GetNewBigCache(ctx, iCacheConfig.bConfig)
		if err != nil {
			log.Println("公司cache 初始化失败")
		}
		OrganizationCache = &ICache{
			Cache:        cache,
			ICacheConfig: iCacheConfig,
		}
	}
}

func InitChannelCache(ctx context.Context, iCacheConfig ICacheConfig) {
	if ChannelCache == nil {
		cache, err := bigcache.New(ctx, bigcache.Config{
			Shards:             16,
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
			log.Println("渠道cache 初始化失败")
		}
		ChannelCache = &ICache{
			Cache:        cache,
			ICacheConfig: iCacheConfig,
		}
	}
}

func ICacheGet(cache *ICache, key string) (value []byte, err error) {
	value, err = cache.Cache.Get(key)
	if err != nil || len(value) == 0 {
		log.Println("key 不存在")
		resMap := map[string]interface{}{}
		resMap, err = cache.ICacheConfig.LoadFunc()
		value, err = json.Marshal(resMap[key])
		if len(value) > 0 {
			err = ICacheSet(cache, key, value)
			if err != nil {
				log.Println("set 失败", err.Error())
			}
		}
		return
	}
	return
}

func ICacheReLoadFunc(ctx context.Context, cache *ICache) {
	startTime := time.Now()
	firstLoad(cache)
	end := time.Since(startTime)
	log.Println(cache.ICacheConfig.CacheName, "cache 加载耗时：", end)
	go func() {
		for {
			select {
			case <-time.After(cache.ICacheConfig.LoadInterval):
				ICacheLoadFunc(cache)
			case <-ctx.Done():
				log.Println("ctx Done")
				break
			}
		}
	}()
	return
}

func firstLoad(cache *ICache) {
	resMap := map[string]interface{}{}
	var err error
	resMap, err = cache.ICacheConfig.LoadFunc()
	if err != nil {
		log.Println("加载失败", err.Error())
	}
	allKey := map[string]struct{}{}
	for k, v := range resMap {
		entry, _ := json.Marshal(v)
		err = ICacheSet(cache, k, entry)
		if err == nil {
			allKey[k] = struct{}{}
		}
		if err != nil {
			log.Println("set失败", err.Error())
		}
	}
	cache.AllKey = allKey
	return
}

func ICacheLoadFunc(cache *ICache) {
	resMap := map[string]interface{}{}
	var err error
	resMap, err = cache.ICacheConfig.LoadFunc()
	if err != nil {
		log.Println("加载失败", err.Error())
	}
	newBigCache := &bigcache.BigCache{}
	newBigCache, err = GetNewBigCache(context.Background(), cache.ICacheConfig.bConfig)
	if newBigCache == nil {
		log.Println("新初始化cache 失败", err)
		return
	}
	allKey := map[string]struct{}{}
	for k, v := range resMap {
		entry, _ := json.Marshal(v)
		err = newBigCache.Set(k, entry)
		if err == nil {
			allKey[k] = struct{}{}
		}
		if err != nil {
			log.Println("set失败", err.Error())
		}
	}
	cache.AllKey = allKey
	cache.Cache = newBigCache
	return
}

func ICacheSet(cache *ICache, key string, entry []byte) error {
	return cache.Cache.Set(key, entry)
}

// GetAllKey key无序
func GetAllKey(cache *ICache) map[string]struct{} {
	return cache.AllKey
}

