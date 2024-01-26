/**
* @Author zdm
* @Date 2023/4/19 16:04
* @Discription
**/
package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var Client *redis.Client
var ReCluster *redis.ClusterClient
var CacheType int

func OnConnect(ctx context.Context, cn *redis.Conn) error {
	if cn != nil {
		fmt.Println("connect ok")
	}
	return nil
}
func InitRedis() {
	if Client == nil {
		fmt.Println("init client")
		Client = redis.NewClient(&redis.Options{
			Addr:      "39.105.142.114:6379",
			Password:  "", // no password set
			DB:        0,  // use default DB
			OnConnect: OnConnect,
		})
	}
}

func InitRedisCluster() {
	if ReCluster == nil {
		ReCluster = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{"39.105.142.114:7000", "39.105.142.114:7001", "39.105.142.114:7002", "39.105.142.114:7003", "39.105.142.114::7004", "39.105.142.114::7005"},
			// To route commands by latency or randomly, enable one of the following.
			//RouteByLatency: true,
			//RouteRandomly: true,
		})
		CacheType = 1
	}
}

// Set key val
func Set(ctx context.Context, key string, val interface{}, duration time.Duration) error {
	if CacheType == 1 {
		return ReCluster.Set(ctx, key, val, duration).Err()
	}
	return Client.Set(ctx, key, val, duration).Err()
}

// Get key
func Get(ctx context.Context, key string) string {
	if CacheType == 1 {
		val, err := ReCluster.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}
		return val
	}
	val, err := Client.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

// HSet key val:map[string][string]
func HSet(ctx context.Context, key string, hKey, val interface{}) error {
	if CacheType == 1 {
		return ReCluster.HSet(ctx, key, hKey, val).Err()
	}
	return Client.HSet(ctx, key, hKey, val).Err()
}

//HGetAll key val:map[string][string]
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	if CacheType == 1 {
		return ReCluster.HGetAll(ctx, key).Result()
	}
	return Client.HGetAll(ctx, key).Result()
}

// LPush key val
func LPush(ctx context.Context, key string, val interface{}) error {
	if CacheType == 1 {
		return ReCluster.LPush(ctx, key, val).Err()
	}
	return Client.LPush(ctx, key, val).Err()
}

func RPop(ctx context.Context, key string) (string, error) {
	if CacheType == 1 {
		return ReCluster.RPop(ctx, key).Result()
	}
	return ReCluster.RPop(ctx, key).Result()
}

func BRPop(ctx context.Context, key string) ([]string, error) {
	if CacheType == 1 {
		return ReCluster.BRPop(ctx, 6*time.Second, key).Result()
	}
	return ReCluster.BRPop(ctx, 6*time.Second, key).Result()
}

func SetNx(ctx context.Context, key string, val string, expTime time.Duration) error {
	if CacheType == 1 {
		return ReCluster.SetNX(ctx, key, val, expTime).Err()
	}
	return Client.SetNX(ctx, key, val, expTime).Err()
}

func Del(ctx context.Context, key string) error {
	if CacheType == 1 {
		return ReCluster.Del(ctx, key).Err()
	}
	return Client.Del(ctx, key).Err()
}

// SAdd add a new member to a set
func SAdd(ctx context.Context, key string, val interface{}) error {
	return ReCluster.SAdd(ctx, key, val).Err()
}

// SIsMember 判断一个某个值在set 中是否存在
func SIsMember(ctx context.Context, key string, val string) (bool, error) {
	if CacheType == 1 {
		return ReCluster.SIsMember(ctx, key, val).Result()
	}
	return Client.SIsMember(ctx, key, val).Result()
}

// SInter 判断key 交集的
func SInter(ctx context.Context, key1, key2 string) ([]string, error) {
	if CacheType == 1 {
		return ReCluster.SInter(ctx, key1, key2).Result()
	}
	return Client.SInter(ctx, key1, key2).Result()
}

// SCard 统计set 的长度
func SCard(ctx context.Context, key string) (int64, error) {
	if CacheType == 1 {
		return ReCluster.SCard(ctx, key).Result()
	}
	return Client.SCard(ctx, key).Result()
}

// SRem 删除对应key set的值
func SRem(ctx context.Context, key string, val ...interface{}) (int64, error) {
	if CacheType == 1 {
		return ReCluster.SRem(ctx, key, val).Result()
	}
	return Client.SRem(ctx, key, val).Result()
}

// ZAdd adds a new member and associated sorted set,If the member already exists,the score is updated
func ZAdd(ctx context.Context, key string, member redis.Z) (int64, error) {
	if CacheType == 1 {
		return ReCluster.ZAdd(ctx, key, member).Result()
	}
	return Client.ZAdd(ctx, key, member).Result()
}

// ZRange returns members of a sorted set,sorted within a given range
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if CacheType == 1 {
		return ReCluster.ZRange(ctx, key, start, stop).Result()
	}
	return Client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeByScore Returns the sorted set cardinality (number of elements) of the sorted set stored at key
func ZRangeByScore(ctx context.Context, key string, by *redis.ZRangeBy) ([]string, error) {
	if CacheType == 1 {
		return ReCluster.ZRangeByScore(ctx, key, by).Result()
	}
	return Client.ZRangeByScore(ctx, key, by).Result()
}

// ZRank returns the rank of the provided member,assuming the sorted is in ascending order
func ZRank(ctx context.Context, key, member string) (int64, error) {
	if CacheType == 1 {
		return ReCluster.ZRank(ctx, key, member).Result()
	}
	return Client.ZRank(ctx, key, member).Result()
}

// ZRevRank returns the rank of the provided member,assuming the sorted set is in descending order.
func ZRevRank(ctx context.Context, key, member string) (int64, error) {
	if CacheType == 1 {
		return ReCluster.ZRevRank(ctx, key, member).Result()
	}
	return Client.ZRevRank(ctx, key, member).Result()
}
