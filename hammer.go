package main

import (
	_ "encoding/json"
	_ "fmt"
	_ "io/ioutil"
	_ "os"
	// "github.com/gin-gonic/gin"
	// "github.com/midoks/hammer/engine"
	"github.com/midoks/hammer/indexer"
)

func main() {
	indexer.ReadConf("conf", func(cf *indexer.ConfFile) {
		indexer.Run(cf)
	})

	// // 初始化
	// searcher.Init(engine.EngineInitOptions{SegmenterDictionaries: "../data/dictionary.txt"})
	// defer searcher.Close()

	// r := gin.Default()
	// r.GET("/pings", func(c *gin.Context) {
	// 	searcher.IndexDocument(1, engine.DocumentIndexData{Content: "此次百度收购将成中国互联网最大并购"}, true)
	// 	c.JSON(200, gin.H{"message": "hammer"})
	// })

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "hammer"})
	// })
	// r.Run(":10000")
}
