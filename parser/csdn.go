package parser

import (
	"CrawlerX/duck"
	"CrawlerX/mod"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/url"
	"runtime/debug"
	"strings"
)

type Csdn func()

// 基础解析器
func (csdn *Csdn) Parse(base *url.URL, reader io.Reader, paths []string) duck.Result {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recover from panic %v", r)
			debug.PrintStack()
		}
	}()
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Printf("read html document error %v", err)
		return nil
	}
	// find all <a> tage
	var nexts []string
	doc.Find("a").Each(func(i int, selection *goquery.Selection) {
		next, ok := selection.Attr("href")
		if ok {
			if next != "" {
				info, err := url.Parse(next)
				if err == nil {
					if info.Host == "" || info.Scheme == "" {
						info.Host = base.Host
						info.Scheme = base.Scheme
					}
					if info.Host == base.Host {
						//name := strings.Trim(selection.Text(), "")
						seed := info.String()
						//log.Printf("find new seed name:%s url:%s", name, seed)
						if isNext(info, paths) {
							nexts = append(nexts, seed)
						}
					}
				}
			}
		}
	})
	return csdn.getResult(doc, nexts)
}

func (csdn *Csdn) getResult(doc *goquery.Document, next []string) duck.Result {
	result := &mod.Topic{}
	// title
	var description, keywords string
	doc.Find("meta").Each(func(i int, selection *goquery.Selection) {
		metaName, ok := selection.Attr("name")
		if ok {
			switch metaName {
			case "description":
				description, _ = selection.Attr("content")
			case "keywords":
				keywords, _ = selection.Attr("content")
			}
		}
	})
	_ = description
	_ = keywords
	doc.Find("link").Each(func(i int, selection *goquery.Selection) {
		rel, ok := selection.Attr("rel")
		if ok {
			if strings.Contains(strings.ToLower(rel), "ico") {
			}
		}
	})
	result.SetNext(next)
	result.V = result
	return result
}
