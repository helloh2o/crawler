package core

import (
	"CrawlerX/duck"
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

func newWorker(id int, rate int) *worker {
	w := &worker{id: id}
	w.req = NewReq("")
	w.in = make(chan string)
	w.out = make(chan duck.Result)
	w.rate = time.NewTicker(time.Millisecond * time.Duration(rate))
	return w
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
				result := parser.Parse(req, reader, PushSeed)
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
