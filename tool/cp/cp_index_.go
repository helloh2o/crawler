package main

import (
	"CrawlerX/db"
	"CrawlerX/mod"
	"CrawlerX/search"
	"encoding/json"
	"flag"
	"github.com/olivere/elastic/v7"
	"log"
)

var (
	from = flag.String("from", "", "copy index from")
	to   = flag.String("to", "", "copy index to")
)

func main() {
	flag.Parse()
	if *from == "" || *to == "" {
		panic("Please write the correct from index and to index.")
	}
	if err := db.InitESClient("http://127.0.0.1:9200", false); err != nil { //docker
		panic(err)
	}
	search.Scroll(*from, elastic.NewMatchAllQuery(), 500, 65535, func(current []*elastic.SearchHit) {
		//copy
		count := 0
		for _, h := range current {
			var page mod.PageInfo
			if err := json.Unmarshal(h.Source, &page); err != nil {
				log.Printf("json unmarshal error %s", err.Error())
			} else {
				if err = db.SavePage(*to, &page); err != nil {
					log.Printf("Copy page to index %s error %s", *to, err.Error())
				} else {
					count++
				}
			}
		}
		log.Printf("Copied data %d", count)
	})
}
