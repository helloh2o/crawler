package main

import (
	"CrawlerX/config"
	"CrawlerX/core"
	"CrawlerX/db"
	"CrawlerX/mod"
	"CrawlerX/out"
	"log"
)

func main() {
	config.Init("./config.yaml")
	if err := db.OpenMySql(config.Instance.MySqlUrl, config.Instance.MySqlMaxIdle, config.Instance.MySqlMaxOpen, config.Instance.ShowSQL, &mod.PageInfo{}); err != nil {
		log.Printf("mysql error %v", err)
	}
	if err := db.InitESClient(config.Instance.ESNode, false); err != nil { //docker
		log.Printf("InitESClient client error %v", err)
	} else {
		if err = db.CreateIndex(config.Instance.EsIndex); err != nil {
			panic(err)
		}
	}
	// 启动爬虫模块
	core.StartEngine(config.Instance.Sites)
	// 订阅输出结果
	out.Save(core.SubResult())
	select {}
}
