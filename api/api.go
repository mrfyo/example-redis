package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitHandler(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")
	router.StaticFile("/favicon.ico", "./assets/favicon.ico")
	homeHandler(router)

	userHandler(router)
	articleHandler(router)
}

func userHandler(router *gin.Engine) {
	router.GET("/api/users", listUserHandler)
	router.POST("/api/users", addUserHandler)
	router.DELETE("/api/users/:id", deleteUserHandler)
}

func articleHandler(r *gin.Engine) {
	r.GET("/api/articles", listArticleHandler)
	r.POST("/api/articles", addArticleHandler)
	r.DELETE("/api/articles/:id", deleteArticleHandler)
	r.GET("/api/articles/top", topListArticleHandler)
	r.POST("/api/articles/vote", voteArticleHandler)
	r.POST("/api/articles/publish", publishArticleHandler)
}

func homeHandler(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
}
