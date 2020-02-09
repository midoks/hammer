package api

import (
	_ "fmt"
	"github.com/gin-gonic/gin"
)

func Run() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hammer"})
	})
	r.Run(":80")
}
