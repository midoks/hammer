package engine

import (
	"fmt"
	"github.com/huichen/murmur"
	"log"
	"runtime"
	// "sync/atomic"
)

type Engine struct {
	initialized bool

	// segmenterChannel chan segmenterRequest
	segmenter EngineSegmenter
}

// 索引器返回结果
type IndexedDocument struct {
	DocId uint64

	// BM25，仅当索引类型为FrequenciesIndex或者LocationsIndex时返回有效值
	BM25 float32

	// 关键词在文档中的紧邻距离，紧邻距离的含义见computeTokenProximity的注释。
	// 仅当索引类型为LocationsIndex时返回有效值。
	TokenProximity int32

	// 紧邻距离计算得到的关键词位置，和Lookup函数输入tokens的长度一样且一一对应。
	// 仅当索引类型为LocationsIndex时返回有效值。
	TokenSnippetLocations []int

	// 关键词在文本中的具体位置。
	// 仅当索引类型为LocationsIndex时返回有效值。
	TokenLocations [][]int
}

// // 评分规则通用接口
// type ScoringCriteria interface {
// 	// 给一个文档评分，文档排序时先用第一个分值比较，如果
// 	// 分值相同则转移到第二个分值，以此类推。
// 	// 返回空切片表明该文档应该从最终排序结果中剔除。
// 	Score(doc IndexedDocument, fields interface{}) []float32
// }

// // 一个简单的评分规则，文档分数为BM25
// type RankByBM25 struct {
// }

// func (rule RankByBM25) Score(doc IndexedDocument, fields interface{}) []float32 {
// 	return []float32{doc.BM25}
// }

// type RankOptions struct {
// 	// 文档的评分规则，值为nil时使用Engine初始化时设定的规则
// 	ScoringCriteria ScoringCriteria

// 	// 默认情况下（ReverseOrder=false）按照分数从大到小排序，否则从小到大排序
// 	ReverseOrder bool

// 	// 从第几条结果开始输出
// 	OutputOffset int

// 	// 最大输出的搜索结果数，为0时无限制
// 	MaxOutputs int
// }

// // 见http://en.wikipedia.org/wiki/Okapi_BM25
// // 默认值见engine_init_options.go
// type BM25Parameters struct {
// 	K1 float32
// 	B  float32
// }

// // 初始化索引器选项
// type IndexerInitOptions struct {
// 	// 索引表的类型，见上面的常数
// 	IndexType int

// 	// 待插入索引表文档 CACHE SIZE
// 	DocCacheSize int

// 	// BM25参数
// 	BM25Parameters *BM25Parameters
// }

// type EngineInitOptions struct {
// 	// 是否使用分词器
// 	// 默认使用，否则在启动阶段跳过SegmenterDictionaries和StopTokenFile设置
// 	// 如果你不需要在引擎内分词，可以将这个选项设为true
// 	// 注意，如果你不用分词器，那么在调用IndexDocument时DocumentIndexData中的Content会被忽略
// 	NotUsingSegmenter bool

// 	// 半角逗号分隔的字典文件，具体用法见
// 	// sego.Segmenter.LoadDictionary函数的注释
// 	SegmenterDictionaries string

// 	// 停用词文件
// 	StopTokenFile string

// 	// 分词器线程数
// 	NumSegmenterThreads int

// 	// 索引器和排序器的shard数目
// 	// 被检索/排序的文档会被均匀分配到各个shard中
// 	NumShards int

// 	// 索引器的信道缓冲长度
// 	IndexerBufferLength int

// 	// 索引器每个shard分配的线程数
// 	NumIndexerThreadsPerShard int

// 	// 排序器的信道缓冲长度
// 	RankerBufferLength int

// 	// 排序器每个shard分配的线程数
// 	NumRankerThreadsPerShard int

// 	// 索引器初始化选项
// 	IndexerInitOptions *IndexerInitOptions

// 	// 默认的搜索选项
// 	DefaultRankOptions *RankOptions

// 	// 是否使用持久数据库，以及数据库文件保存的目录和裂分数目
// 	UsePersistentStorage    bool
// 	PersistentStorageFolder string
// 	PersistentStorageShards int
// }

// func (engine *Engine) Init(options EngineInitOptions) {
// 	// 将线程数设置为CPU数
// 	runtime.GOMAXPROCS(runtime.NumCPU())

// 	// 初始化初始参数
// 	if engine.initialized {
// 		log.Fatal("请勿重复初始化引擎")
// 	}
// 	engine.initialized = true

// 	// if !options.NotUsingSegmenter {
// 	// 	// 载入分词器词典
// 	// 	engine.segmenter.LoadDictionary(options.SegmenterDictionaries)
// 	// 	// 初始化停用词
// 	// 	engine.stopTokens.Init(options.StopTokenFile)
// 	// }

// 	// 初始化索引器和排序器
// 	// for shard := 0; shard < options.NumShards; shard++ {
// 	// 	engine.indexers = append(engine.indexers, core.Indexer{})
// 	// 	engine.indexers[shard].Init(*options.IndexerInitOptions)

// 	// 	engine.rankers = append(engine.rankers, core.Ranker{})
// 	// 	engine.rankers[shard].Init()
// 	// }

// 	// 初始化分词器通道
// 	// engine.segmenterChannel = make(
// 	// 	chan segmenterRequest, options.NumSegmenterThreads)

// 	engine.segmenter.Init()

// 	// 启动分词器
// 	// for iThread := 0; iThread < options.NumSegmenterThreads; iThread++ {

// 	// }

// }

// func (engine *Engine) IndexDocument(docId uint64, data DocumentIndexData, forceUpdate bool) {
// 	if !engine.initialized {
// 		log.Fatal("必须先初始化引擎")
// 	}

// 	hash := murmur.Murmur3([]byte(fmt.Sprint("%d%s", docId, data.Content)))
// 	fmt.Println(hash)
// 	engine.segmenter.channel <- segmenterRequest{docId: docId, hash: hash, data: data, forceUpdate: forceUpdate}
// }

// // 关闭引擎
// func (engine *Engine) Close() {
// }
