package main

import (
	_ "encoding/json"
	_ "fmt"
	"github.com/midoks/hammer/api"
	"github.com/midoks/hammer/configure"
	"github.com/midoks/hammer/indexer"
	"github.com/robfig/cron"
	_ "io/ioutil"
	"log"
	_ "os"
	_ "time"
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

	configure.Read("conf", func(cf *configure.Args) {
		go indexer.Run(cf)
	})

	api.Run()
}
