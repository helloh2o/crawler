package main

import (
	"CrawlerX/config"
	"CrawlerX/core"
	"CrawlerX/db"
	"CrawlerX/out"
	"CrawlerX/parser"
	"log"
)

func main() {
	config.Init("./config.yaml")
	/*if err := db.OpenMySql(config.Instance.MySqlUrl, config.Instance.MySqlMaxIdle, config.Instance.MySqlMaxOpen, config.Instance.ShowSQL, &mod.PageInfo{}); err != nil {
		log.Println(err)
	} else {
		db.MysqlReady = true
	}*/
	if err := db.InitESClient(config.Instance.ESNode, false); err != nil { //docker
		log.Printf("InitESClient client error %v", err)
	} else {
		db.EsReady = true
		if err = db.CreateIndex(config.Instance.EsIndex); err != nil {
			panic(err)
		}
	}
	// 获取结果
	out.Save(core.SubResult())
	// 设置参数
	params := config.Instance.CrawlParma
	core.InitEngine(params.MaxWaitQueueSize, params.MaxWorkers, params.WorkerRate, config.Instance.Parsers)()
	// 为解析器设置站点
	parser.SetSitesMap(config.Instance.SiteMap)
	// 抓取指定网站列表
	for _, site := range config.Instance.Sites {
		core.PushSeed(site.Seed)
	}
	select {}
}
