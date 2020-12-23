package parser

import (
	"CrawlerX/duck"
	"CrawlerX/mod"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/url"
	"regexp"
	"runtime/debug"
	"strings"
	"time"
)

type PageBasicParser func()

// 基础解析器
func (cmm *PageBasicParser) Parse(base *url.URL, reader io.Reader, paths []string) duck.Result {
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
	return cmm.getResult(doc, base, nexts)
}

func isNext(info *url.URL, paths []string) bool {
	if len(paths) == 0 {
		return true
	} else {
		for _, rule := range paths {
			ruleSlice := strings.Split(rule, "/")
			regex := "^"
			if len(ruleSlice) > 0 && ruleSlice[0] == "*" {
				regex = ""
			}
			for _, per := range ruleSlice {
				if per == "" {
					continue
				}
				flag := "/"
				if per == "*" {
					per = ".*"
					flag = ""
				}
				regex += flag + per
			}
			matched, err := regexp.MatchString(regex, info.Path)
			//log.Printf("regex:: %s, path::%s matched %v Host %s", regex, info.Path, matched,info.Host)
			if err != nil {
				continue
			}
			if matched {
				return matched
			}
		}
	}
	return false
}

func (cmm *PageBasicParser) getResult(doc *goquery.Document, base *url.URL, next []string) duck.Result {
	result := &mod.PageInfo{}
	// title
	title := doc.Find("title").Text()
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
	doc.Find("link").Each(func(i int, selection *goquery.Selection) {
		rel, ok := selection.Attr("rel")
		if ok {
			if strings.Contains(strings.ToLower(rel), "ico") {
				if icoHref, ok := selection.Attr("href"); ok {
					result.ICO = icoHref
				}
			}
		}
	})
	result.Domain = base.Host
	result.Title = title
	result.KeyWords = keywords
	result.Description = description
	result.URL = base.String()
	result.CreateAt = time.Now().Unix()
	result.SetNext(next)
	result.V = result
	return result
}
