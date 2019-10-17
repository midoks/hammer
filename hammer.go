package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/midoks/hammer/engine"
)

var (
	// searcher是线程安全的
	searcher = engine.Engine{}
)

func main() {

	// 初始化
	searcher.Init(types.EngineInitOptions{
		SegmenterDictionaries: "../data/dictionary.txt"})
	defer searcher.Close()

	r := gin.Default()
	r.GET("/pings", func(c *gin.Context) {
		fmt.Println("ss")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hammer",
		})
	})
	r.Run(":10000")
}
