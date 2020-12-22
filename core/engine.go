package core

import (
	"CrawlerX/duck"
	"CrawlerX/mod"
)

var (
	output = make(chan duck.Result, 10000)
)

// 订阅输出结果通道
func SubResult() <-chan duck.Result {
	return output
}

// 初始化引擎
func StartEngine(targets []*mod.Site) <-chan duck.Result {
	for i := 0; i < len(targets); i++ {
		// 创建工作组
		runWorkerGroup(i+1, targets[i], output)
	}
	return output
}
