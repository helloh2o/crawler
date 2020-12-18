package mod

import "CrawlerX/duck"

type TopicDocument struct {
	Id              int64
	NodeId          int64
	UserId          int64
	Title           string
	Content         string
	Recommend       bool
	LastCommentTime int64
	Status          int
	ViewCount       int64
	CommentCount    int64
	LikeCount       int64
	CreateTime      int64
}

func (td *TopicDocument) Value() duck.Result { return td }
