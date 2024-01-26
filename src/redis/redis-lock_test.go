package redis

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// redis
func redisNxLock() {
	InitRedisCluster()

	for i := 0; i < 1000; i++ {
		go func() {
			ctx := context.Background()
			err := SetNx(ctx, "key", "key", 0)
			if err == nil {
				dec()
				Del(ctx, "key")
			}
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	}
	time.Sleep(1 * time.Second)
	fmt.Println(count)
}
func TestLock(t *testing.T) {
	redisNxLock()
}

var count = 10000

func dec() {
	count = count - 1
}
