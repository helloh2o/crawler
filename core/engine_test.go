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
		WorkerSize: 100,
		WorkerRate: 50,
		Parser:     new(parser.Jianshu),
	}})
	time.Sleep(time.Minute * 10)
}
