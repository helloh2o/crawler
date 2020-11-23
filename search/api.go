package search

import (
	"CrawlerX/mod"
	"encoding/json"
	irisCtx "github.com/kataras/iris/v12/context"
	"github.com/olivere/elastic/v7"
	"log"
	"strconv"
)

func init() {
	app.Get("/search/{keywords}", func(ctx *irisCtx.Context) {
		keywords := ctx.Params().Get("keywords")
		search(ctx, keywords, getIndex(ctx))
	})
	app.Get("/search", func(ctx *irisCtx.Context) {
		keywords := ctx.FormValue("query")
		search(ctx, keywords, getIndex(ctx))
	})
}
func getIndex(ctx *irisCtx.Context) string {
	index := ctx.FormValue("type")
	if index == "" && len(Config.IndexKey) == 1 {
		for _, v := range Config.IndexKey {
			return v
		}
	}
	return index
}
func search(ctx *irisCtx.Context, keywords string, indexName string) {
	if keywords == "" {
		ctx.Redirect("/", 302)
	}
	pageStr := ctx.FormValue("pn")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	var data []mod.PageInfo
	var query elastic.Query
	// 精确匹配
	//query = elastic.NewTermQuery("title", keywords)
	filedS := Config.IndexQuery[indexName]
	// not found the query filedS
	if len(filedS) == 0 {
		ctx.StatusCode(500)
		return
	}
	var total interface{} = 0
	var hits *elastic.SearchHits
	var pageData []*elastic.SearchHit
	// 循环匹配
	for _, qName := range filedS {
		// 习语匹配查询 + ik_smart 分词匹配
		query = elastic.NewMatchPhraseQuery(qName, keywords).Analyzer("ik_smart")
		// 模糊匹配
		blurQuery := elastic.NewMatchQuery(qName, keywords).Analyzer("ik_smart")
		// 默认得分排序获取最优结果
		hits, pageData = Scroll(indexName, query, 10, page, nil)
		if hits != nil && hits.TotalHits != nil {
			total = hits.TotalHits.Value
		}
		// 没有找到结果，进行模糊匹配, 再找
		if total == 0 {
			query = blurQuery
			hits, pageData = Scroll(indexName, query, 10, page, nil)
			if hits != nil && hits.TotalHits != nil {
				total = hits.TotalHits.Value
			}
			ctx.ViewData("tips", "根据您的关键字没有找到较为精确的结果！相似的查找，")
		} else {
			// 匹配到数据
			break
		}
	}
	ctx.ViewData("hits", total)
	//log.Printf("Result data %v", data)
	// 首页权重查询
	wSize := 0
	/*if page == 1 {
		var weightHits []*elastic.SearchHit
		search := db.GetClientByIndex(db.EsIndex).Search(db.EsIndex).Size(10).From(0).Sort("weight", false)
		resp, err := search.Query(query).Do(context.Background())
		if err == nil {
			weightHits = resp.Hits.Hits
			total := len(weightHits)
			// 权重结果暂时取前 <=5个
			if total <= 5 {
				weightHits = weightHits[:total]
			} else {
				weightHits = weightHits[:5]
			}
		} else {
			log.Printf("Get weight result error %v", err.Error())
		}
		// 合并最优结果
		wSize = len(weightHits)
		if wSize > 0 && len(pageData) >= wSize {
			pageData = append(weightHits, pageData[:len(pageData)-wSize]...)
		}
		log.Printf("Got weight result  %d", wSize)
	}*/
	// 渲染view
	for i, h := range pageData {
		var d mod.PageInfo
		err := json.Unmarshal(h.Source, &d)
		if err != nil {
			log.Println(err)
		} else {
			if wSize-1 >= i {
				d.Title = "===权重命中===" + d.Title
			}
			data = append(data, d)
		}
	}
	if ctx.Request().URL.Path == "/search" {
		view(keywords, data, ctx)
	} else {
		_, err := ctx.JSON(data)
		if err != nil {
			ctx.StatusCode(500)
		}
	}
}

func view(keywords string, data []mod.PageInfo, ctx *irisCtx.Context) {
	ctx.ViewData("data", data)
	ctx.ViewData("title", "搜索 - "+keywords)
	if err := ctx.View("result.html"); err != nil {
		ctx.StatusCode(502)
	}
}
