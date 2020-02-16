package configure

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type ArgsConn struct {
	Localhost string
	Port      int
	User      string
	Password  string
	Db        string
	Charset   string
}

type Args struct {
	Path     string
	AppName  string
	Type     string
	Conn     ArgsConn
	Pk       string
	Query    string
	Step     int
	Start    int64
	Interval string
}

/** 替换注释 */
func ReplaceConfComment(s string) string {
	re := regexp.MustCompile(`(//.*)`)
	s = re.ReplaceAllString(s, "")

	re2 := regexp.MustCompile(`(?s)(/*(.*?)*/)`)
	s = re2.ReplaceAllString(s, "")
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
