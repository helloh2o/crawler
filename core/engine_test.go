package core

import (
	"CrawlerX/duck"
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
	InitEngine(1000, 2, 200, ps)()
	PushSeed("https://mlog.club/")
	select {}
}
