package core

import (
	"CrawlerX/duck"
	"CrawlerX/mod"
	"CrawlerX/parser"
	"testing"
)

func TestPushSeed(t *testing.T) {
	ps := make(map[string]duck.Parser)
	ps["default"] = new(parser.PageBasicParser)
	go func() {
		out := SubResult()
		for {
			<-out
		}
	}()
	InitEngine([]mod.Site{mod.Site{
		Seed:           "",
		Weight:         0,
		Paths:          nil,
		ExpirationDays: 0,
		ParserName:     "",
		WorkerSize:     0,
		WorkerRate:     0,
	}}, ps)()
	select {}
}
