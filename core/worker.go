package core

import (
	"CrawlerX/config"
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
	site *mod.Site
}

func runWorkerGroup(id int, site *mod.Site, output chan duck.Result) {
	input := make(chan string, 10000*site.WorkerSize)
	req := NewHttpClient(config.Instance.Proxy)
	rate := time.NewTicker(time.Millisecond * time.Duration(site.WorkerRate))
	for i := 0; i < site.WorkerSize; i++ {
		w := &worker{id: id}
		w.req = req
		w.in = input
		w.out = output
		w.rate = rate
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
			history.Store(target, true)
			log.Printf("worker %d get target %s", w.id, target)
			<-w.rate.C
			// 请求 (TODO 分布式)
			go w.DoParse(target) // 同步
		}
	}()
}

func (w *worker) DoParse(target string) {
	w.req.DoReq("GET", target, func(req *url.URL, reader io.Reader) {
		result := w.site.Parser.Parse(req, reader, w.site.Paths)
		if result != nil {
			if result.Value() != nil {
				output <- result
			}
			if len(result.GetNext()) > 0 {
				for _, task := range result.GetNext() {
					if _, ok := history.Load(task); ok {
						return
					}
					select {
					case w.in <- task:
						log.Printf("add next task %s", task)
					default:
						log.Printf("worker is busy.")
					}
				}
			}
		}
	})
}

func (w *worker) Run() {
	w.Consume()
}
