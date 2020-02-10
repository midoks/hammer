package configure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ArgsConn struct {
	Localhost string
	Port      int
	User      string
	Pwd       string
	Db        string
	Charset   string
}

type Args struct {
	Type  string
	Name  string
	Pk    string
	Conn  ArgsConn
	Sql   string
	Step  int
	Start int
}

func Read(path string, call func(conf *Args)) {
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		panic("conf is error!")
		return
	}
	for _, fi := range rd {
		j := fmt.Sprintf("%s/%s/data.json", path, fi.Name())

		f, err := os.Open(j)
		c, err := ioutil.ReadAll(f)

		if err != nil {
			fmt.Printf("path:conf/%s is error!", fi.Name())
			continue
		}

		a := &Args{}
		err = json.Unmarshal([]byte(c), &a)
		if err == nil {
			call(a)
		} else {
			fmt.Println("path:conf/%s json has error!", fi.Name())
		}
	}
}
