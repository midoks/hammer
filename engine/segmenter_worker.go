package engine

import (
	"fmt"
	// "log"
)

// 文档的一个关键词
type TokenData struct {
	// 关键词的字符串
	Text string

	// 关键词的首字节在文档中出现的位置
	Locations []int
}

type DocumentIndexData struct {
	// 文档全文（必须是UTF-8格式），用于生成待索引的关键词
	Content string

	// 文档的关键词
	// 当Content不为空的时候，优先从Content中分词得到关键词。
	// Tokens存在的意义在于绕过悟空内置的分词器，在引擎外部
	// 进行分词和预处理。
	Tokens []TokenData

	// 文档标签（必须是UTF-8格式），比如文档的类别属性等，这些标签并不出现在文档文本中
	Labels []string

	// 文档的评分字段，可以接纳任何类型的结构体
	Fields interface{}
}

type segmenterRequest struct {
	docId       uint64
	hash        uint32
	data        DocumentIndexData
	forceUpdate bool
}

func (engine *Engine) segmenterWorker() {
	for {
		request := <-engine.segmenterChannel
		fmt.Println("segmenterWorker start")
		fmt.Println(request)
		fmt.Println("segmenterWorker end")
	}
}

// 关闭引擎
func (engine *Engine) Close() {
}
