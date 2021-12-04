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

	user := User{
		Name: "Jack",
	}

	if err := CreateUser(&user); err != nil {
		fmt.Println(err)
		return
	}

	art, err := GetUser(user.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", art)

}
