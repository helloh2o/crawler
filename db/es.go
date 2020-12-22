package db

import (
	"CrawlerX/mod"
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"reflect"
	"strconv"
)

var (
	EsClient *elastic.Client
)

// 初始化ES
func InitESClient(url string, sniff bool) (err error) {
	ctx := context.Background()
	EsClient, err = elastic.NewClient(elastic.SetSniff(sniff), elastic.SetURL(url))
	if err != nil {
		return
	}
	info, code, err := EsClient.Ping(url).Do(ctx)
	if err != nil {
		return
	}
	log.Printf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)
	return
}

func SavePage(index string, data interface{}) (err error) {
	if EsClient == nil {
		return
	}
	ctx := context.Background()
	p, ok := data.(*mod.PageInfo)
	if ok {
		id := strconv.Itoa(int(p.Id))
		_, err = EsClient.Index().Index(index).Id(id).BodyJson(&p).Do(ctx)
		return err
	}
	return errors.New(fmt.Sprintf("Unknown type %v", reflect.TypeOf(data)))
}

func DeleteIndex(index string) {
	if EsClient == nil {
		return
	}
	ctx := context.Background()
	if _, err := EsClient.DeleteIndex(index).Do(ctx); err != nil {
		log.Printf("Delete index %s failed.", index)
	}
}
func CreateIndex(indexName string) error {
	if EsClient == nil {
		return errors.New("es not open")
	}
	ctx := context.Background()
	exists, err := EsClient.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}
	// create
	if !exists {
		_, err = EsClient.CreateIndex(indexName).BodyString(getDefaultSettings()).Do(ctx)
		if err != nil {
			return err
		}
	} else {
		return nil
	}
	log.Printf("created index %s", indexName)
	return err
}

// 索引的数据量、可能的并发数要求，还有es本身的限制，一个分片最多能索引20亿条数据,
// 默认分片3个，shard *3 = lucene index*3， 副本一个， 默认分词器ik_max_word,
func getDefaultSettings() string {
	return `{
				  "settings": {
					"index": {
  					  "analysis.analyzer.default.type": "ik_max_word",
					  "number_of_shards": "3",
					  "number_of_replicas": "1"
					}
				  }
				}`
}
