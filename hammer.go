package main

import (
	_ "encoding/json"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	_ "io/ioutil"
	"log"
	_ "os"
	_ "time"
	// "github.com/midoks/hammer/engine"
	"github.com/midoks/hammer/indexer"
)

func cronInit() {
	// go func() {
	crontab := cron.New()

	i := 0
	log.Println("Run Test AddFunc")
	crontab.AddFunc("* * * * * ?", func() {
		log.Println("ddd")
		i++
		log.Println("cron running:", i)

	})

	// crontab.AddFunc("@every 1s", print15)
	crontab.Start()
	defer crontab.Stop()
	// }()
}

func print5() {
	log.Println("Run 5s cron")
}

func print10() {
	log.Println("Run 10s cron")
}

func print15() {
	log.Println("Run 15s cron")
}

func main() {

	cronInit()

	// t1 := time.NewTimer(time.Second * 2)
	// for {
	// 	select {
	// 	case <-t1.C:
	// 		t1.Reset(time.Second * 2)
	// 		print10()
	// 	}
	// }

	// select {}
	// defer c.Stop()

	indexer.ReadConf("conf", func(cf *indexer.ConfFile) {
		indexer.Run(cf)
	})

	r := gin.Default()
	// r.GET("/pings", func(c *gin.Context) {
	// 	searcher.IndexDocument(1, engine.DocumentIndexData{Content: "此次百度收购将成中国互联网最大并购"}, true)
	// 	c.JSON(200, gin.H{"message": "hammer"})
	// })

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hammer"})
	})
	r.Run(":80")

}
