package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	redisDB *redis.Client
	ctx     = context.Background()
)

func main() {
	CtrlCWatcher()

	redisDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "wait_123456",
		DB:       1,
	})

	router := gin.Default()

	userHandler(router)

	router.Run(":8080")
}

func userHandler(router *gin.Engine) {
	router.GET("/users", ListUser)
	router.POST("/users", AddUser)
	router.DELETE("/users/:id", DeleteUser)
}

func CtrlCWatcher() {
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
