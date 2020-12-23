package mod

type Model struct {
	Id int64 `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
}

type Topic struct {
	Model
	NodeId          int64       `gorm:"not null;index:idx_node_id;" json:"nodeId" form:"nodeId"`                         // 节点编号
	UserId          int64       `gorm:"not null;index:idx_topic_user_id;" json:"userId" form:"userId"`                   // 用户
	Title           string      `gorm:"size:128" json:"title" form:"title"`                                              // 标题
	Content         string      `gorm:"type:longtext" json:"content" form:"content"`                                     // 内容
	Recommend       bool        `gorm:"not null;index:idx_recommend" json:"recommend" form:"recommend"`                  // 是否推荐
	ViewCount       int64       `gorm:"not null" json:"viewCount" form:"viewCount"`                                      // 查看数量
	CommentCount    int64       `gorm:"not null" json:"commentCount" form:"commentCount"`                                // 跟帖数量
	LikeCount       int64       `gorm:"not null" json:"likeCount" form:"likeCount"`                                      // 点赞数量
	Status          int         `gorm:"index:idx_topic_status;" json:"status" form:"status"`                             // 状态：0：正常、1：删除
	LastCommentTime int64       `gorm:"index:idx_topic_last_comment_time" json:"lastCommentTime" form:"lastCommentTime"` // 最后回复时间
	CreateTime      int64       `gorm:"index:idx_topic_create_time" json:"createTime" form:"createTime"`                 // 创建时间
	ExtraData       string      `gorm:"type:text" json:"extraData" form:"extraData"`                                     // 扩展数据
	next            []string    `gorm:"-" json:"next"`
	V               interface{} `gorm:"-" json:"v"`
}

func (tp *Topic) SetNext(tasks []string) {
	tp.next = tasks
}
func (tp *Topic) GetNext() []string {
	return tp.next
}
func (tp *Topic) Value() interface{} {
	return tp.V
}
