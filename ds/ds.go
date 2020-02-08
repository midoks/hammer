package ds

import (
	"fmt"
)

type DateSourceIf interface {
	Import()
	Task()
	GetData() (map[int]map[string]string, error)
}

func Factory(name string) DateSourceIf {
	switch name {
	case "mysql":
		ds := &DataSourceMySQL{}
		ds.Init()
		return ds
	default:
		panic("No such animal")
	}
}

func Run() {

	fmt.Println("ds!")
}
