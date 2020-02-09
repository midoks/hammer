package storage

import (
	"fmt"
)

type StorageIf interface {
	Add(map[string]string)
}

func Factory(name string) StorageIf {
	switch name {
	case "lucene":
		sl := &StorageLucene{}
		return sl
	default:
		panic("No such animal")
	}
}

func Run() {
	fmt.Println("storage!")
}
