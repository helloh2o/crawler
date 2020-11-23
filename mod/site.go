package mod

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
}
