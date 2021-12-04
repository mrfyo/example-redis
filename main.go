package main

import (
	"context"
	"fmt"

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

	if err := redisDB.Set(ctx, "hi", "No", 0).Err(); err != nil {
		fmt.Println(err)
		return
	}

	res := redisDB.Get(ctx, "hi")
	if err := res.Err(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.Val())
}
