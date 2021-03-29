package parser

import (
	"CrawlerX/duck"
	"CrawlerX/mod"
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"runtime/debug"
	"strconv"
	"time"
)

//https://i.news.qq.com/trpc.qqnews_web.kv_srv.kv_srv_http_proxy/list?sub_srv_id=fashion&srv_id=pc&offset=0&limit=10&strategy=1&ext={%22pool%22:[%22top%22],%22is_filter%22:2,%22check_type%22:true}
type QQNews struct{}
type JsonRet struct {
	Ret  int
	Msg  string
	Data map[string][]Article
}

type Article struct {
	Title        string `json:"title"`
	Url          string `json:"url"`
	Thumb        string `json:"img"`
	CategoryName string `json:"category_name"`
	CategoryCN   string `json:"category_cn"`
	Tags         []Tag  `json:"tags"`
}

type Tag struct {
	Name string `json:"tag_word"`
}

/**
	## 单行的标题
	**粗体**
	`console.log('行内代码')`
	```js\n code \n``` 标记代码块
	[内容](链接)
	![文字说明](图片链接)
**/

// 基础解析器
func (qq *QQNews) Parse(base *url.URL, reader io.Reader, paths []string) duck.Result {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recover from panic %v", r)
			debug.PrintStack()
		}
	}()
	data, err := ioutil.ReadAll(reader)
	var nexts []string
	// json 结果
	if string(data[0]) == "{" {
		var ret JsonRet
		if err = json.Unmarshal(data, &ret); err != nil {
			log.Printf("unmarshal qq JsonRet error %v", err)
			return nil
		}
		offset, _ := strconv.Atoi(base.Query().Get("offset"))
		limit, _ := strconv.Atoi(base.Query().Get("limit"))
		base.Query().Set("offset", strconv.Itoa(offset+limit))
		for _, page := range ret.Data["list"] {
			query := "?category_name=" + page.CategoryName
			query += "&category_cn=" + page.CategoryCN
			tags := ""
			for _, t := range page.Tags {
				if tags != "" {
					tags += "_"
				}
				tags += t.Name
			}
			query += "&tags=" + tags
			nexts = append(nexts, page.Url+query)
		}
		nexts = append(nexts, base.String())
		result := &mod.Topic{
			CreateTime: time.Now().UnixNano() / int64(time.Millisecond),
		}
		result.V = nil
		result.SetNext(nexts)
		return result
	} else {
		var err error
		data, err = GbkToUtf8(data)
		if err != nil {
			log.Printf("utf8 to gbk error %v", err)
			return nil
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
		//log.Println(doc.Html())
		if err != nil {
			log.Printf("read html document error %v", err)
			return nil
		}
		result := &mod.Topic{
			CreateTime: time.Now().UnixNano() / int64(time.Millisecond),
		}
		result.LastCommentTime = result.CreateTime
		result.Title = doc.Find("h1").Text()
		result.NodeId = 1
		result.UserId = 1
		result.Recommend = true
		content, err := doc.Find(".content-article").Html()
		if err != nil {
			log.Printf("find content-article html error %v", err)
			return nil
		}
		result.Content = Convert(content)
		if result.Content == "" {
			result.V = nil
		} else {
			result.V = result
		}
		return result
	}
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
