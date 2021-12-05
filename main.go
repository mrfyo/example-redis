package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/mrfyo/example-redis/api"
	"github.com/mrfyo/example-redis/model"
)

var (
	redisDB *redis.Client
	ctx     = context.Background()
)

func main() {
	InterruptWatcher()

	redisDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "wait_123456",
		DB:       1,
	})

	router := gin.Default()

	model.InitModel(redisDB, ctx)
	api.InitHandler(router)

	router.Run(":8080")
}

func InterruptWatcher() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGALRM)
	go func() {
		for sig := range c {
			log.Printf("captured %v, stopping profiler and exiting..", sig)

			redisDB.Close()
			os.Exit(0)
		}
	}()
}
