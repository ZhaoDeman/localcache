package redis

import (
	"fmt"
	redisbloom "github.com/RedisBloom/redisbloom-go"
)

// redis set 判断key是否存在，可能会占用内存
// 采用布隆过滤器

func test() {
	var client = redisbloom.NewClient("localhost:6379", "nohelp", nil)

	// BF.ADD mytest item
	_, err := client.Add("mytest", "myItem")
	if err != nil {
		fmt.Println("Error:", err)
	}

	exists, err := client.Exists("mytest", "myItem")
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("myItem exists in mytest: ", exists)
}
