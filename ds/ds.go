package ds

import (
	"encoding/json"
	"fmt"
	"github.com/midoks/hammer/configure"
	"io/ioutil"
	"os"
	"time"
)

var (
	DS_TYPE_MYSQL = "mysql"
)

type DateSourceIf interface {
	Import()
	Task()
	DeltaData()
	DeleteData()
}

type SaveStatus struct {
	PK              int64  `json:"pk"`
	LastUpdatedTime string `json:"last_updated_time"`
}

func (ss *SaveStatus) Read(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	c, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(c), ss)
	if err != nil {
		return err
	}
	return nil
}

func (ss *SaveStatus) Save(pk int64, path string) error {
	ctime := time.Now().Format("2006-01-02 15:04:05")
	ss.PK = pk
	ss.LastUpdatedTime = ctime

	b, err := json.Marshal(ss)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, b, 0777)
	if err != nil {
		return err
	}

	return nil
}

func OpenDS(conf *configure.Args) DateSourceIf {
	switch conf.Type {
	case DS_TYPE_MYSQL:
		ds := &DataSourceMySQL{}
		ds.Init(conf)
		return ds
	default:
		ds := &DataSourceMySQL{}
		ds.Init(conf)
		return ds
	}
}

func Run() {
	fmt.Println("ds!")
}
