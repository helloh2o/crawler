package core

import (
	"CrawlerX/duck"
	"CrawlerX/mod"
	"io"
	"log"
	"net/url"
	"sync"
	"time"
)

type worker struct {
	id      int
	in      chan string
	out     chan duck.Result
	rate    *time.Ticker
	req     *Req
	site    *mod.Site
	history sync.Map
}

func runWorkerGroup(id int, site *mod.Site, output chan duck.Result) {
	input := make(chan string, 10000*site.WorkerSize)
	req := NewReq("")
	for i := 0; i < site.WorkerSize; i++ {
		w := &worker{id: id}
		w.req = req
		w.in = input
		w.out = output
		w.rate = time.NewTicker(time.Millisecond * time.Duration(site.WorkerRate))
		w.site = site
		go w.Run()
	}
	// seed
	input <- site.Seed
}

func (w *worker) Consume() {
	go func() {
		for {
			target := <-w.in
			w.history.Store(target, true)
			log.Printf("worker %d get target %s", w.id, target)
			<-w.rate.C
			// 请求
			go w.req.Crawl(target, func(req *url.URL, reader io.Reader) {
				result := w.site.Parser.Parse(req, reader, w.site.Paths, func(task string) {
					if _, ok := w.history.Load(task); ok {
						return
					}
					select {
					case w.in <- task:
					default:
					}
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
