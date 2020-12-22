package mod

import "CrawlerX/duck"

type Site struct {
	// 种子页面
	Seed string `yaml:"seed"`
	// 站点权重
	Weight int `yaml:"weight"`
	// 路径包含
	Paths []string `yaml:"paths"`
	// 过期时间, 天
	ExpirationDays int `yaml:"days"`
	// 解析器
	ParserName string `yaml:"parser"`
	// 工人数量
	WorkerSize int `yaml:"worker_size"`
	// 工人频率
	WorkerRate int `yaml:"worker_rate"`
	// 解析器
	Parser duck.Parser `yaml:"-"`
}
