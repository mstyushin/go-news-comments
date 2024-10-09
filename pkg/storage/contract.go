package storage

import "context"

type Comment struct {
	ID        int    `json:"id"`
	ArticleID int    `json:"article_id"`
	ParentID  int    `json:"parent_id"`
	Author    string `json:"author"`
	Text      string `json:"text"`
	PubTime   int64  `json:"pub_time"`
}

type Storage interface {
	AddComment(ctx context.Context, comment Comment) (int, error)
	GetCommentsByArticleID(ctx context.Context, articleID int) ([]Comment, error)
	DeleteComment(ctx context.Context, commentID int) error
}
