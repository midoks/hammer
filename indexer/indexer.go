package indexer

import (
	"encoding/json"
	"fmt"
	"github.com/midoks/hammer/ds"
	"github.com/midoks/hammer/storage"
	_ "github.com/robfig/cron"
	"io/ioutil"
	"os"
)

type ConfFileConn struct {
	Localhost string
	Port      int
	User      string
	Pwd       string
	Db        string
}

type ConfFile struct {
	Type  string
	Name  string
	Pk    string
	Conn  ConfFileConn
	Sql   string
	Step  int
	Start int
}

func ReadConf(path string, call func(conf *ConfFile)) {
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		panic("conf is error!")
		return
	}
	for _, fi := range rd {
		confFile := fmt.Sprintf("%s/%s/data.json", path, fi.Name())

		filePtr, err := os.Open(confFile)
		c, err := ioutil.ReadAll(filePtr)

		if err != nil {
			fmt.Printf("path:conf/%s is error!", fi.Name())
			continue
		}

		cf := &ConfFile{}
		err = json.Unmarshal([]byte(c), &cf)
		if err == nil {
			call(cf)
		} else {
			fmt.Println("path:conf/%s json has error!", fi.Name())
		}
	}
}

func Init() {
}

func Run(cf *ConfFile) {

	dsObj := ds.Factory(cf.Type)

	// go dsObj.Import()
	// go dsObj.Task()

	storage.Run()

	// for {
	// 	d := dsObj.DataChan
	// 	fmt.Println(d)
	// }

	v, err := dsObj.GetData()

	fmt.Println(v, err)
}
