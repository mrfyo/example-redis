package main

import (
	"context"
	"log"
	"net/http"
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
	InterruptWatcher()

	redisDB = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "wait_123456",
		DB:       1,
	})

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")
	HomeHandler(router)
	UserHandler(router)
	ArticleHandler(router)
	router.Run(":8080")
}

func HomeHandler(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
}

func UserHandler(router *gin.Engine) {
	router.GET("/api/users", ListUserHandler)
	router.POST("/api/users", AddUserHandler)
	router.DELETE("/api/users/:id", DeleteUserHandler)
}

func ArticleHandler(r *gin.Engine) {
	r.GET("/api/articles", ListArticleHandler)
	r.POST("/api/articles", AddArticleHandler)
	r.DELETE("/api/articles/:id", DeleteArticleHandler)
	r.GET("/api/articles/top", TopListArticleHandler)
	r.POST("/api/articles/vote", VoteArticleHandler)
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
