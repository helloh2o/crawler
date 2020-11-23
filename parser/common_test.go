package parser

import (
	"CrawlerX/mod"
	"log"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

func TestIsNext(t *testing.T) {
	info, _ := url.Parse("https://blog.csdn.net/34324/article/details/78749419")
	site := new(mod.Site)
	site.Paths = []string{"*/article/details"}
	for _, rule := range site.Paths {
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
		log.Printf("regex:: %s, path::%s matched %v", regex, info.Path, matched)
		if err != nil {
			continue
		}
		log.Printf("%v", matched)
	}
}
