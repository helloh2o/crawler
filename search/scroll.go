package search

import (
	"CrawlerX/db"
	"context"
	"github.com/olivere/elastic/v7"
	"log"
)

// Scroll 100 items per search
func Scroll(index string, query elastic.Query, pSize int, pNum int, deepFunc func(current []*elastic.SearchHit)) (hits *elastic.SearchHits, page []*elastic.SearchHit) {
	deep, max := 1, 1
	res, err := db.EsClient.Scroll().Index(index).
		Query(query).
		Scroll("5m").
		Size(pSize).
		Do(context.Background())
	if err != nil {
		log.Printf("Scroll error %v", err)
	} else {
		if deepFunc != nil {
			// 每个深度回调执行
			deepFunc(res.Hits.Hits)
		}
		max = int(res.Hits.TotalHits.Value) / pSize
		log.Printf("Total hits %d, max page %d", res.Hits.TotalHits.Value, max)
		//checkResult(res.Hits.Hits)
		for {
			// 深度页
			if pNum == deep || deep > max {
				return res.Hits, res.Hits.Hits
			}
			// delay for the deep search
			//time.Sleep(time.Millisecond * 10)
			res, err = db.EsClient.Scroll("1m").ScrollId(res.ScrollId).Do(context.TODO())
			if err == nil {
				if len(res.Hits.Hits) <= 0 {
					break
				}
				deep++
				if deepFunc != nil {
					// 每个深度回调执行
					deepFunc(res.Hits.Hits)
				}
			} else {
				if err.Error() != "EOF" {
					log.Printf("Scroll error %v", err)
				}
				break
			}
		}
	}
	return
}
