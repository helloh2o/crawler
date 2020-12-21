package core

import (
	"CrawlerX/duck"
	"CrawlerX/mod"
	"log"
	"os"
)

var (
	workers []*worker
	sites   []mod.Site
	stop    = make(chan struct{})
	output  = make(chan duck.Result)
	parsers map[string]duck.Parser
)

// 订阅输出结果通道
func SubResult() <-chan duck.Result {
	return output
}

// 初始化引擎
func start() {
	// collect result
	go func() {
		for _, w := range workers {
			select {
			case r := <-w.out:
				select {
				case output <- r:
				default:
					log.Printf("//******// Unhandle Result From Worker %d :: %+v", w.id, r)
				}
			default:
				continue
			}
		}
	}()
	createWkGroups()
	<-stop
	os.Exit(0)
}

// 创建工作组
func createWkGroups() {
	for i := 0; i < len(sites); i++ {
		runWorkerGroup(i+1, sites[i], output)
	}
}

// 设置解析器
func InitEngine(targets []mod.Site, ps map[string]duck.Parser) func() {
	parsers = ps
	sites = targets
	return start
}
