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
	GetData() (map[int]map[string]string, error)
}

type SaveStatus struct {
	ID          int64  `json:"id"`
	CurrentTime string `json:"current_time"`
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

func (ss *SaveStatus) Save(id int64, path string) error {
	ctime := time.Now().Format("2006-01-02 15:04:05")
	ss.ID = id
	ss.CurrentTime = ctime

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

func Factory(conf *configure.Args) DateSourceIf {
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
