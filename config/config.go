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
	MySqlUrl     string     `yaml:"MySqlUrl"`
	MySqlMaxIdle int        `yaml:"MySqlMaxIdle"`
	MySqlMaxOpen int        `yaml:"MySqlMaxOpen"`
	ShowSQL      bool       `yaml:"ShowSQL"`
	CrawlParma   crawlParma `yaml:"CrawlParma"`
	// ES node
	ESNode  string `yaml:"ESNode"`
	EsIndex string `yaml:"EsIndex"`
	// 抓取列表
	Sites   []mod.Site             `yaml:"Sites"`
	SiteMap map[string]mod.Site    `yaml:"-"`
	Parsers map[string]duck.Parser `yaml:"-"`
}
type crawlParma struct {
	MaxWaitQueueSize int `yaml:"MaxWaitQueueSize"`
	MaxWorkers       int `yaml:"MaxWorkers"`
	WorkerRate       int `yaml:"WorkerRate"`
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
		Instance.SiteMap[info.Host] = site
		registerParseByName(site)
	}
}

// 注册配置解析器
func registerParseByName(site mod.Site) {
	switch site.ParserName {
	case "":
		Instance.Parsers[""] = new(parser.PageBasicParser)
	case "default":
		Instance.Parsers["default"] = new(parser.PageBasicParser)
	default:
		Instance.Parsers["default"] = new(parser.PageBasicParser)
	}
}
