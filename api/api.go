package api

import (
	_ "fmt"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "hammer"})
}

func Run() {
	r := gin.Default()
	r.GET("/ping", Ping)
	r.GET("/so", So)
	r.Run(":80")
}
