package main

import (
	"context"

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

	defer redisDB.Close()

	// article := &Article{
	// 	ID:     0,
	// 	Title:  "Hi Go",
	// 	Link:   "http://127.0.0.1/Go",
	// 	Poster: "wait",
	// 	Time:   time.Now().Unix(),
	// 	Votes:  0,
	// }

	// CreateArticle(article)

	// user := &User{
	// 	Name: "Jack",
	// }

	// CreateUser(user)

	// VoteArticle(user, article)

	// article, _ := GetArticle(1)
	// user, _ := GetUser(2)

	// VoteArticle(user, article)
}
