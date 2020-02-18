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

	name := ""
	for _, v := range data {
		name += fmt.Sprintf("%s|==|", v)
	}

	d := document.NewDocument()                                              //创建doucument
	d.Add(document.NewTextFieldFromString("text", name, document.STORE_YES)) //添加域信息
	fmt.Println("d.Fields():", d.Fields())
	writer.AddDocument(d.Fields())
	defer writer.Close()
}

func (sl *StorageLucene) Search(key string) {
	util.SetDefaultInfoStream(util.NewPrintStreamInfoStream(os.Stdout))
	index.DefaultSimilarity = func() index.Similarity {
		return search.NewDefaultSimilarity() //评分器
	}
	directory, _ := store.OpenFSDirectory("conf/test/data") //创建索引的目录
	// analyzer := std.NewStandardAnalyzer()                   //使用标准分词器，貌似只有这一个分词器其他的没见到
	// conf := index.NewIndexWriterConfig(util.VERSION_LATEST, analyzer) //Indexwriter的配置器

	reader, _ := index.OpenDirectoryReader(directory) //打开reader
	searcher := search.NewIndexSearcher(reader)       //创建searcher

	q := search.NewTermQuery(index.NewTerm("text", key)) //termquery  目前只发现有termquery 和 boolean Query 其他的span profix phrase 这些应该是没有 因为用的term所以是产于分词的分词器是单字切分，所以只用一个雨字来搜索
	res, _ := searcher.Search(q, nil, 1000)              //result search 中传入query filter 和返回的条数
	fmt.Printf("Found %v hit(s).\n", res.TotalHits)
	for _, hit := range res.ScoreDocs {
		fmt.Printf("Doc %v score: %v\n", hit.Doc, hit.Score)
		doc, _ := reader.Document(hit.Doc)
		fmt.Printf("text -> %v\n", doc.Get("text"))
	}

}
