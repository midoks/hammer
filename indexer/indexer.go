package indexer

import (
	"fmt"
	"github.com/midoks/hammer/configure"
	"github.com/midoks/hammer/cron"
	"github.com/midoks/hammer/ds"
	"log"
)

func Run(cf *configure.Args) {

	ods := ds.OpenDS(cf)

	go ods.Import()
	go ods.Task()

	cron.Add("@every 3s", func() {
		log.Println("indexr cron! 3s")
		r, _ := ods.GetData()
		fmt.Println(r)
	})

}
