package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisDB *redis.Client
	ctx     = context.Background()
)

func main() {

	redisDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "wait_123456",
		DB:       1,
	})

	key := "article:123"

	err := redisDB.HSet(ctx, key, map[string]interface{}{
		"title":  "Redis Reference",
		"link":   "http://doc.redisfans.com/",
		"poster": "1",
		"time":   time.Now().Unix(),
		"votes":  32,
	}).Err()

	if err != nil {
		fmt.Println(err)
		return
	}

	res := redisDB.HGetAll(ctx, key)
	if err := res.Err(); err != nil {
		fmt.Println(err)
		return 
	}

	fmt.Println(res.Val())
	
}
