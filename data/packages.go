package data

import (
	"./providers/commoncrawl"
	"./processors/nlp/knowledge"
)

type Packages struct {
	Commoncrawl commoncrawl.F
	Knowledge knowledge.F
}