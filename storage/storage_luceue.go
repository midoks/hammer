package storage

import (
	"fmt"
	std "github.com/balzaczyy/golucene/analysis/standard"
	_ "github.com/balzaczyy/golucene/core/analysis/tokenattributes"
	_ "github.com/balzaczyy/golucene/core/codec/lucene410"
	"github.com/balzaczyy/golucene/core/document"
	"github.com/balzaczyy/golucene/core/index"
	"github.com/balzaczyy/golucene/core/search"
	"github.com/balzaczyy/golucene/core/store"
	"github.com/balzaczyy/golucene/core/util"
	// "log"
	"os"
	// "strings"
)

type StorageLucene struct {
}

func (sl *StorageLucene) Add(data map[string]string) {
	util.SetDefaultInfoStream(util.NewPrintStreamInfoStream(os.Stdout))
	index.DefaultSimilarity = func() index.Similarity {
		return search.NewDefaultSimilarity() //评分器
	}

	directory, _ := store.OpenFSDirectory("conf/test/data")           //创建索引的目录
	analyzer := std.NewStandardAnalyzer()                             //使用标准分词器，貌似只有这一个分词器其他的没见到
	conf := index.NewIndexWriterConfig(util.VERSION_LATEST, analyzer) //Indexwriter的配置器
	writer, _ := index.NewIndexWriter(directory, conf)                //index writer

	for k, v := range data {
		fmt.Println(k, v)
		// d := document.NewDocument()
		//创建doucument
		// d.Add(document.NewTextFieldFromString("text", v, document.STORE_YES)) //添加域信息
		// writer.AddDocument(d.Fields())
	}

	d := document.NewDocument()                                                      //创建doucument
	d.Add(document.NewTextFieldFromString("text", data["name"], document.STORE_YES)) //添加域信息
	fmt.Println("d.Fields():", d.Fields())
	writer.AddDocument(d.Fields())
	defer writer.Close()
}
