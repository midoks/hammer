package storage

import (
	"fmt"
)

var (
	ENGINE_TYPE_LUCENE = "lucene"
	ENGINE_TYPE_WUKONG = "wukong"
)

type StorageIf interface {
	Add(map[string]string)
}

func OpenStorage(name string) StorageIf {
	switch name {
	case ENGINE_TYPE_LUCENE:
		sl := &StorageLucene{}
		return sl
	case ENGINE_TYPE_WUKONG:
		sw := &StorageWukong{}
		return sw
	default:
		panic("No such animal")
	}
}

func Run() {
	fmt.Println("storage!")
}
