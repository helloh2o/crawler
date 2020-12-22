package core

import (
	"CrawlerX/mod"
	"CrawlerX/parser"
	"log"
	"testing"
	"time"
)

func TestPushSeed(t *testing.T) {
	go func() {
		out := SubResult()
		for {
			result := <-out
			log.Printf("Got ==>\n %v", result.Value())
		}
	}()
	StartEngine([]*mod.Site{&mod.Site{
		Seed:       "https://www.jianshu.com/",
		Paths:      []string{"/p/*"},
		WorkerSize: 4,
		WorkerRate: 200,
		Parser:     new(parser.PageBasicParser),
	}})
	time.Sleep(time.Minute * 10)
}
