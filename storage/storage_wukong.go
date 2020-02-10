package storage

import (
	"github.com/huichen/wukong/engine"
	"github.com/huichen/wukong/types"
	"log"
)

type StorageWukong struct {
}

var (
	// searcher是协程安全的
	searcher = engine.Engine{}
)

func (sl *StorageWukong) Add(data map[string]string) {

	// 初始化
	searcher.Init(types.EngineInitOptions{SegmenterDictionaries: "test/dictionary.txt"})
	defer searcher.Close()

	// 将文档加入索引，docId 从1开始
	searcher.IndexDocument(1, types.DocumentIndexData{Content: "此次百度收购将成中国互联网最大并购"}, false)

	// 等待索引刷新完毕
	searcher.FlushIndex()
}

func (sl *StorageWukong) Search(key string) {
	// 初始化
	searcher.Init(types.EngineInitOptions{
		SegmenterDictionaries: "test/dictionary.txt"})
	defer searcher.Close()

	// 搜索输出格式见types.SearchResponse结构体
	log.Print(searcher.Search(types.SearchRequest{Text: "百度中国"}))

}
