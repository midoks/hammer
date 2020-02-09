package api

import (
	"fmt"
	_ "github.com/balzaczyy/golucene/analysis/standard"
	_ "github.com/balzaczyy/golucene/core/analysis/tokenattributes"
	_ "github.com/balzaczyy/golucene/core/codec/lucene410"
	"github.com/balzaczyy/golucene/core/index"
	"github.com/balzaczyy/golucene/core/search"
	"github.com/balzaczyy/golucene/core/store"
	"github.com/balzaczyy/golucene/core/util"
	// "log"
	"os"
	// "strings"
	"github.com/gin-gonic/gin"
)

func So(c *gin.Context) {
	util.SetDefaultInfoStream(util.NewPrintStreamInfoStream(os.Stdout))
	index.DefaultSimilarity = func() index.Similarity {
		return search.NewDefaultSimilarity() //评分器
	}
	directory, _ := store.OpenFSDirectory("conf/test/data") //创建索引的目录
	// analyzer := std.NewStandardAnalyzer()                   //使用标准分词器，貌似只有这一个分词器其他的没见到
	// conf := index.NewIndexWriterConfig(util.VERSION_LATEST, analyzer) //Indexwriter的配置器

	reader, _ := index.OpenDirectoryReader(directory) //打开reader
	searcher := search.NewIndexSearcher(reader)       //创建searcher

	key := c.Query("q")
	fmt.Println(key)

	q := search.NewTermQuery(index.NewTerm("name", key)) //termquery  目前只发现有termquery 和 boolean Query 其他的span profix phrase 这些应该是没有 因为用的term所以是产于分词的分词器是单字切分，所以只用一个雨字来搜索
	res, _ := searcher.Search(q, nil, 1000)              //result search 中传入query filter 和返回的条数
	fmt.Printf("Found %v hit(s).\n", res.TotalHits)
	for _, hit := range res.ScoreDocs {
		fmt.Printf("Doc %v score: %v\n", hit.Doc, hit.Score)
		doc, _ := reader.Document(hit.Doc)
		fmt.Printf("text -> %v\n", doc.Get("text"))
	}

	c.JSON(200, gin.H{"message": "hammer"})
}
