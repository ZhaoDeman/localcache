# Redis List

## Redis Cli
redis客户端连接
```go
func InitRedis() {
	if Client == nil {
		Client = redis.NewClient(&redis.Options{
			Addr:      "127.0.0.1:6379",
			Password:  "", // no password set
			DB:        0,  // use default DB
			OnConnect: nil,
		})
	}
}
```

## List
利用List实现异步、解耦、削峰
> LPUSH、BRPOP 实现先进先出的队列

```go
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
```