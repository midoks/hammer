package ds

import (
	"fmt"
)

type DateSourceIf interface {
	Import() bool
	GetData() string
}

func Factory(name string) DateSourceIf {
	switch name {
	case "mysql":
		return &DataSourceMySQL{}
	default:
		panic("No such animal")
	}
}

func Run() {

	fmt.Println("ds!")
}
