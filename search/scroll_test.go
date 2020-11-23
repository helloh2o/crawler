package search

import (
	"CrawlerX/db"
	"github.com/olivere/elastic/v7"
	"log"
	"testing"
	"time"
)

func TestExpiration(t *testing.T) {
	err := db.InitESClient("http://localhost:9200", false)
	if err != nil {
		panic(err)
	}
	boolQry := elastic.NewBoolQuery()
	boolQry.Must(elastic.NewRangeQuery("expiration").Gt(time.Now().Unix()))
	hits, result := Scroll("pages", boolQry, 100, 2000, func(current []*elastic.SearchHit) {
		/*log.Printf("hits lenght %d", len(current))
		for _, h := range current {
			log.Printf("Id -> %s", h.Id)
		}*/
	})
	log.Printf("Total %d", hits.TotalHits.Value)
	for _, h := range result {
		log.Printf("Id -> %s", h.Id)
	}
}
