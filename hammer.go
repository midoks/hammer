package main

import (
	"github.com/midoks/hammer/api"
	"github.com/midoks/hammer/configure"
	"github.com/midoks/hammer/indexer"
	"runtime"
)

func initSetting() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// initSetting()

	configure.Read("conf", func(cf *configure.Args) {
		go indexer.Run(cf)
	})

	api.Run()
}
