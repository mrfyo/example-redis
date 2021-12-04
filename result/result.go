package result

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Fail(c *gin.Context, code int, message string) {
	if code == 0 {
		code += 100
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
	})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": data,
	})
}
