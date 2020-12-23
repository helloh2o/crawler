package mod

// 网页结构体
type PageInfo struct {
	Id int64 `gorm:"primarykey" json:"id"`
	// 域名
	Domain string `gorm:"index:idx_domain;COMMENT:'域名'" json:"domain"`
	// 标题
	Title string `gorm:"COMMENT:'标题'" json:"title"`
	// 描述
	Description string `gorm:"COMMENT:'描述'" json:"description"`
	// 关键字
	KeyWords string `gorm:"COMMENT:'关键字'" json:"key_words"`
	// 地址
	URL string `gorm:"unique;index:idx_url;COMMENT:'地址'" json:"url"`
	// 点击次数
	Clicks int64 `gorm:"COMMENT:'点击次数'" json:"clicks"`
	// 时间
	CreateAt int64 `gorm:"COMMENT:'时间'" json:"create_at"`
	// 搜索权重
	Weight int `gorm:"COMMENT:'搜索权重'" json:"weight"`
	// 过期时间
	Expiration int64 `gorm:"COMMENT:'过期时间'" json:"expiration"`
	// ico
	ICO  string      `gorm:"COMMENT:'网站小图标'" json:"ico"`
	next []string    `gorm:"-" json:"next"`
	V    interface{} `gorm:"-" json:"v"`
}

func (pg *PageInfo) SetNext(tasks []string) {
	pg.next = tasks
}
func (pg *PageInfo) GetNext() []string {
	return pg.next
}
func (pg *PageInfo) Value() interface{} {
	return pg.V
}
