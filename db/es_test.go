package db

import (
	"CrawlerX/mod"
	"log"
	"testing"
)

var _ = `
{
    "settings" : {
        "index" : {
            "analysis.analyzer.default.type": "ik_max_word",
            "number_of_shards" : 5,
            "number_of_replicas" : 1
        }
    }
}
`

func TestCreateES(t *testing.T) {
	if err := InitESClient("http://127.0.0.1:9200", false); err != nil { //docker
		panic(err)
	}
}

func TestSave(t *testing.T) {
	index := "p7"
	if err := InitESClient("http://127.0.0.1:9200", false); err != nil { //docker
		panic(err)
	}
	for i := 0; i < 10; i++ {
		result := mod.PageInfo{Id: int64(i)}
		if err := SavePage(index, &result); err != nil {
			log.Printf("ES save page err %v", err)
		}
	}
}

func TestCreateIndex(t *testing.T) {
	indexs := []string{"index_topic", "index_user", "index_favorite_user"}
	if err := InitESClient("http://192.168.1.150:9200/", false); err != nil { //docker
		log.Printf("InitESClient ES client error %v", err)
	}
	for _, idx := range indexs {
		if err := CreateIndex(idx); err != nil {
			panic(err)
		}
	}
}

func TestDeleteIndex(t *testing.T) {
	index := "index_user"
	if err := InitESClient("http://127.0.0.1:9200", false); err != nil { //docker
		log.Printf("InitESClient ES client error %v", err)
	}
	DeleteIndex(index)
}
