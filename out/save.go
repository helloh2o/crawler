package out

import (
	"CrawlerX/config"
	"CrawlerX/db"
	"CrawlerX/duck"
	"context"
	"log"
)

func Save(out <-chan duck.Result) {
	go func() {
		var count int
		for {
			result := <-out
			if db.DB() != nil {
				if err := db.DB().Save(result).Error; err != nil {
					log.Printf("Mysql save error %v", err)
				}
			}
			indexName := config.Instance.EsIndex
			if db.EsClient != nil {
				if err := db.SavePage(indexName, result); err != nil {
					log.Printf("ES save page err %v", err)
				} else {
					count++
				}
			}
			log.Printf("Save Result Count:: %d", count)
			// 每1000条flush to disk
			if count%1000 == 0 && db.EsClient != nil {
				_, err := db.EsClient.Flush().Index(indexName).Do(context.Background())
				if err != nil {
					log.Printf("es flush error %v", err)
				} else {
					log.Println("es flush data ... ")
				}
			}
		}
	}()
}
