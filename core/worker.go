package core

import (
	"CrawlerX/duck"
	"CrawlerX/mod"
	"io"
	"log"
	"net/url"
	"time"
)

type worker struct {
	id   int
	in   chan string
	out  chan duck.Result
	rate *time.Ticker
	req  *Req
}

func runWorkerGroup(id int, site mod.Site, output chan duck.Result) {
	input := make(chan string)
	for i := 0; i < site.WorkerSize; i++ {
		w := &worker{id: id}
		w.req = NewReq("")
		w.in = input
		w.out = output
		w.rate = time.NewTicker(time.Millisecond * time.Duration(site.WorkerRate))
		go w.Run()
	}
	// seed
	input <- site.Seed
}

func (w *worker) Consume() {
	go func() {
		for {
			target := <-w.in
			log.Printf("worker %d get target %s", w.id, target)
			<-w.rate.C
			// 请求
			go w.req.Crawl(target, func(req *url.URL, reader io.Reader) {
				p, ok := parsers[req.Host]
				var parser duck.Parser
				if ok {
					// 特殊站点解析器
					parser = p
				} else {
					// 默认解析器
					parser = parsers["default"]
				}
				result := parser.Parse(req, reader, func(task string) {
					// 自产自销
					w.in <- task
				})
				if result != nil {
					output <- result
				}
			})
		}
	}()
}

func (w *worker) Run() {
	w.Consume()
}
