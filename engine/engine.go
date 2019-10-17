package engine

import (
	"fmt"
)

type Engine struct {
	segmenterChannel chan segmenterRequest
}
