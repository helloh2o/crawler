package out

import (
	"CrawlerX/config"
	"CrawlerX/db"
	"CrawlerX/duck"
	"CrawlerX/mod"
	"context"
	"log"
	"time"
)

func Save(out <-chan duck.Result) {
	go func() {
		var count int
		for {
			result := <-out
			saved := false
			if db.MysqlReady {
				if err := db.DB().Save(result).Error; err != nil {
					log.Printf("Mysql save error %v", err)
				} else {
					saved = true
				}
			}
			indexName := config.Instance.EsIndex
			if db.EsReady {
				// save to ES
				if !db.MysqlReady {
					p, ok := result.(*mod.PageInfo)
					if ok {
						p.Id = time.Now().Unix() + int64(count+1)
					}
				}
				if err := db.SavePage(indexName, result); err != nil {
					log.Printf("ES save page err %v", err)
				} else {
					saved = true
				}
			}
			if !saved {
				log.Printf("Not Saved the Result :: %v", result)
			} else {
				count++
				log.Printf("Save Result Count:: %d", count)
			}
			// 每1000条flush to disk
			if count%1000 == 0 && db.EsReady {
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
