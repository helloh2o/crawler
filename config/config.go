package config

import (
	"CrawlerX/duck"
	"CrawlerX/mod"
	"CrawlerX/parser"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/url"
)

type Config struct {
	MySqlUrl     string `yaml:"MySqlUrl"`
	MySqlMaxIdle int    `yaml:"MySqlMaxIdle"`
	MySqlMaxOpen int    `yaml:"MySqlMaxOpen"`
	ShowSQL      bool   `yaml:"ShowSQL"`
	// ES node
	ESNode  string `yaml:"ESNode"`
	EsIndex string `yaml:"EsIndex"`
	// 抓取列表
	Sites   []*mod.Site            `yaml:"Sites"`
	SiteMap map[string]mod.Site    `yaml:"-"`
	Parsers map[string]duck.Parser `yaml:"-"`
}

var (
	Instance *Config
)

func Init(filename string) {
	Instance = &Config{}
	if yamlFile, err := ioutil.ReadFile(filename); err != nil {
		log.Fatal(err)
	} else if err = yaml.Unmarshal(yamlFile, Instance); err != nil {
		log.Fatal(err)
	}
	Instance.SiteMap = make(map[string]mod.Site)
	Instance.Parsers = make(map[string]duck.Parser)
	for _, site := range Instance.Sites {
		info, err := url.Parse(site.Seed)
		if err != nil || info.Host == "" {
			log.Fatalf("站点种子配置错误 %v", err)
		}
		registerParseByName(site)
	}
}

// 注册配置解析器
func registerParseByName(site *mod.Site) {
	switch site.ParserName {
	case "csdn":
		site.Parser = new(parser.Csdn)
	default:
		site.Parser = new(parser.PageBasicParser)
	}
}
