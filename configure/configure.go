package configure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type ArgsConn struct {
	Localhost string `json:"localhost"`
	Port      int    `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Db        string `json:"db"`
	Charset   string `json:"charset"`
}

type Args struct {
	Path       string
	AppName    string
	Type       string   `json:"type"`
	Conn       ArgsConn `json:"conn"`
	Pk         string   `json:"pk"`
	Query      string   `json:"query"`
	Step       int      `json:"step"`
	Start      int64    `json:"start"`
	Interval   string   `json:"interval"`
	DeltaQuery string   `json:"delta_query"`
}

/** 替换注释 */
func ReplaceConfComment(s string) string {
	singleLineRe := regexp.MustCompile(`(//.*)`)
	s = singleLineRe.ReplaceAllString(s, "")

	multLine := regexp.MustCompile(`(?s)(/*(.*)*/)`)
	s = multLine.ReplaceAllString(s, "")
	return s
}

/* 监控配置文件 */
func Watcher() {

}

/* 读取配置文件 */
func Read(path string, call func(conf *Args)) {
	flist, err := ioutil.ReadDir(path)
	if err != nil {
		panic("conf is error!")
		return
	}
	for _, fi := range flist {
		j := fmt.Sprintf("%s/%s/data.json", path, fi.Name())

		f, err := os.Open(j)
		c, err := ioutil.ReadAll(f)

		if err != nil {
			fmt.Printf("path:conf/%s is error!", fi.Name())
			continue
		}

		b := ReplaceConfComment(string(c))

		a := &Args{}
		err = json.Unmarshal([]byte(b), &a)
		if err == nil {
			a.Path = path
			a.AppName = fi.Name()
			call(a)
		} else {
			fmt.Println("path:conf/%s json has error!", fi.Name())
		}
	}
}
