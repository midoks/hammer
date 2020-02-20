package configure

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	Path             string
	AppName          string
	Type             string   `json:"type"`
	Conn             ArgsConn `json:"conn"`
	Pk               string   `json:"pk"`
	Query            string   `json:"query"`
	Step             int      `json:"step"`
	Start            int64    `json:"start"`
	Interval         string   `json:"interval"`
	DeltaQuery       string   `json:"delta_query"`
	DeltaImportQuery string   `json:"delta_import_query"`
	DeletedQuery     string   `json:"deleted_query"`
}

/** 替换注释 */
func ReplaceConfComment(s string) string {
	singleLineRe := regexp.MustCompile(`(//.*)`)
	s = singleLineRe.ReplaceAllString(s, "")

	multLine := regexp.MustCompile(`(?s)(/*(.*)*/)`)
	s = multLine.ReplaceAllString(s, "")
	return s
}

func IsListenConf(cPath string, root string) bool {
	vPath := strings.Replace(cPath, root, "", -1)
	vPath = strings.Trim(vPath, "/")
	vArr := strings.SplitN(vPath, "/", 2)

	vArrLen := len(vArr)
	if vArrLen == 2 && vArr[1] == "data.json" {
		return true
	} else if vArrLen == 1 {
		return true
	}
	return false
}

func ListenConf(root string) {

	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("start:", err)
	}
	defer watch.Close()

	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if IsListenConf(path, root) {
			err := watch.Add(path)
			if err != nil {
				log.Println("filepath data.json file err:", err)
			}
		}
		return nil
	})

	err = watch.Add(root)
	if err != nil {
		log.Println("root err:", err)
		return
	}

	go func() {
		for {
			select {
			case ev := <-watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件 : ", ev.Name)
						if IsListenConf(ev.Name, root) {
							err := watch.Add(ev.Name)
							if err != nil {
								log.Println("ev.Op&fsnotify.Create == fsnotify.Create 创建文件 err:", err)
							}
						}
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("写入文件 : ", ev.Name)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件 : ", ev.Name)
						watch.Remove(ev.Name)
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						log.Println("重命名文件 : ", ev.Name)
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						log.Println("修改权限 : ", ev.Name)
					}
				}
			case err := <-watch.Errors:
				{
					if err != nil {
						log.Println("error: ", err)
					}
					return
				}
			}
		}
	}()
	select {}
}

/* 监控配置文件 */
func Watcher(path string) {
	go ListenConf(path)
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
