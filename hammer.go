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
	"runtime"
	"time"
)

func cronInit() {
	nyc, _ := time.LoadLocation("Asia/Shanghai")
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(nyc))

	log.Println("Run Test AddFunc")
	crontab.AddFunc("@every 3s", func() {
		log.Println("ddd")
	})

	// crontab.AddFunc("@every 1s", print15)
	crontab.Start()
	// defer crontab.Stop()
}

func initSetting() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	// initSetting()
	cronInit()

	configure.Read("conf", func(cf *configure.Args) {
		go indexer.Run(cf)
	})

	api.Run()
}
