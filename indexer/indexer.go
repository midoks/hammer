package indexer

import (
	"fmt"
	"github.com/midoks/hammer/configure"
	"github.com/midoks/hammer/ds"
	_ "github.com/midoks/hammer/storage"
	_ "github.com/robfig/cron"
)

func Run(cf *configure.Args) {

	dsObj := ds.Factory(cf)

	go dsObj.Import()
	go dsObj.Task()

	v, err := dsObj.GetData()

	fmt.Println(v, err)
}
