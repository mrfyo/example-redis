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

	article := Article{
		ID:     0,
		Title:  "Hi Go",
		Link:   "http://127.0.0.1/Go",
		Poster: "wait",
		Time:   time.Now().Unix(),
		Votes:  46,
	}

	if err := CreateArticle(&article); err != nil {
		fmt.Println(err)
		return
	}

	art, err := GetArticle(article.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", art)

}
