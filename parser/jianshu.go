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
	"time"
)

type Jianshu struct{}

/**
	## 单行的标题
	**粗体**
	`console.log('行内代码')`
	```js\n code \n``` 标记代码块
	[内容](链接)
	![文字说明](图片链接)
**/

// 基础解析器
func (js *Jianshu) Parse(base *url.URL, reader io.Reader, paths []string) duck.Result {
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
	return js.getResult(doc, nexts)
}

func (js *Jianshu) getResult(doc *goquery.Document, next []string) duck.Result {
	result := &mod.Topic{
		CreateTime: time.Now().UnixNano() / int64(time.Millisecond),
	}
	result.LastCommentTime = result.CreateTime
	title := doc.Find("title").Text()
	result.Title = strings.Replace(title, "- 简书", "", -1)
	result.NodeId = 1
	result.UserId = 1
	doc.Find("div").Each(func(i int, selection *goquery.Selection) {
		v, ok := selection.Attr("role")
		if ok && v == "main" {
			var err error
			result.Content, err = selection.Html()
			idx := strings.Index(result.Content, "推荐阅读")
			if idx != -1 {
				result.Content = result.Content[:idx]
			}
			if err != nil {
				return
			} else {
				result.Content = Convert(result.Content)
				result.Content = strings.Replace(result.Content, "data-original-src=", "src=", -1)
			}
		}
	})
	if result.Content == "" {
		result.V = nil
	} else {
		result.V = result
	}
	result.SetNext(next)
	return result
}
