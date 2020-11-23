package main

import (
	"CrawlerX/db"
	"CrawlerX/search"
	sv "CrawlerX/search"
	"flag"
)

var conf = flag.String("conf", "./config.yaml", "the config file")

// 记录，
// 1.es数据删除，go-my-es不会同步老数据到ES，需要删除bin记录文件，重新同步数据
// 2.go-my-es默认没有索引，外部创建不能设置mapping否则会冲突，只设置默认参数，分配，备份，默认分词器
func main() {
	flag.Parse()
	search.Init(*conf)
	if err := db.InitESClient(search.Config.EsNode, false); err != nil {
		panic(err)
	}
	for idx, _ := range search.Config.IndexQuery {
		if err := db.CreateIndex(idx); err != nil {
			panic(err)
		}
	}
	sv.RunWeb(search.Config.Addr)
}
