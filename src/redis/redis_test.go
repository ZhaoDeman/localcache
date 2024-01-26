/**
* @Author zdm
* @Date 2023/4/19 20:50
* @Discription
**/
package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	InitRedisCluster()
	ctx := context.Background()
	err := Set(ctx, "zdm", "赵德满", 1*time.Hour)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(Get(ctx, "zdm"))
}

func TestHash(t *testing.T) {
	InitRedis()
	session := map[string]string{"name": "zdm", "company": "sm", "age": "24"}
	ctx := context.Background()
	for k, v := range session {
		err := HSet(ctx, "user:session", k, v)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	res, err := HGetAll(ctx, "user:session")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res)
}

func TestList(t *testing.T) {
	InitRedisCluster()
	ctx := context.Background()
	// insert Gz100 队列中
	_ = LPush(ctx, "audit:rule", "Gz100")
	// insert Gz101 队列中
	_ = LPush(ctx, "audit:rule", "Gz101")
	// insert Gz102 队列中
	_ = LPush(ctx, "audit:rule", "Gz102")

	go func() {
		for i := 0; i < 10; i++ {
			_ = LPush(ctx, "audit:rule", fmt.Sprintf("%s%d", "Gz11", i))
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			_ = LPush(ctx, "audit:rule", fmt.Sprintf("%s%d", "Gz21", i))
			time.Sleep(10 * time.Millisecond)
		}
	}()

	go func() {
		for {
			val, err := BRPop(ctx, "audit:rule")
			if err != nil {
				fmt.Println("err--->",err.Error())
			}
			fmt.Println("协程1：", val)
			//处理业务
			time.Sleep(1 * time.Second)
			// update val 已处理
		}
	}()
	go func() {
		for {
			val, err := BRPop(ctx, "audit:rule")
			if err != nil {
				fmt.Println("err--->",err.Error())
			}
			fmt.Println("协程2：", val)
			time.Sleep(1 * time.Second)
			// update val 已处理
		}
	}()
	time.Sleep(30*time.Second)
	_ = LPush(ctx,"audit:rule","Gz200")
	time.Sleep(100 * time.Second)
}

func TestSet(t *testing.T) {
	key := "user"
	InitRedisCluster()
	ctx := context.Background()
	err := SAdd(ctx,key,"zhaodeman")
	if err != nil {
		fmt.Println("add fail--->",err.Error())
	}
	ok ,err := SIsMember(ctx,key,"zhaodeman")
	if ok {
		fmt.Println("zhaodeman存在")
	}else{
		fmt.Println("zhaodeman不存在")
	}
	ok ,err = SIsMember(ctx,key,"zdm")
	if ok {
		fmt.Println("zdm 存在")
	}else{
		fmt.Println("zdm不存在")
	}
}

func TestZSet(t *testing.T) {
	key := "user:play"
	InitRedisCluster()
	ctx := context.Background()
	index,err := ZAdd(ctx,key,redis.Z{
		Score: 100,
		Member: "zdm",
	})
	fmt.Println(index,"--->",err)
	index,err = ZAdd(ctx,key,redis.Z{
		Score: 90,
		Member: "qmy",
	})
	fmt.Println(index,"--->",err)

	index,err = ZAdd(ctx,key,redis.Z{
		Score: 80,
		Member: "spx",
	})
	fmt.Println(index,"---->",err)

	index,err = ZAdd(ctx,key,redis.Z{
		Score: 110,
		Member: "sjx",
	})
	fmt.Println(index,"--->",err)

	index,err = ZAdd(ctx,key,redis.Z{
		Score: 109,
		Member: "ly",
	})
	fmt.Println(index,"--->",err)

	index,err = ZAdd(ctx,key,redis.Z{
		Score: 109,
		Member: "zyk",
	})
	fmt.Println("ZRang 使用")
	res,err := ZRange(ctx,key,0,2)
	fmt.Println("从小到大的3个--->",res)
	res,err = ZRange(ctx,key,0,10)
	fmt.Println("从小到大的所有内容--->",res)
	rank,err := ZRank(ctx,key,"sjx")
	fmt.Println("查询sjx的从小到大的等级是多少",rank)
	rank,err = ZRank(ctx,key,"zdm")
	fmt.Println("查询zdm的从小到大的等级是多少",rank)

	rank,err = ZRevRank(ctx,key,"zdm")
	fmt.Println("查询zdm的从大到小的等级是多少",rank)
}
