package pg

import (
	"context"
	"errors"
	"log"

	"github.com/mstyushin/go-news-comments/pkg/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

var _ storage.Storage = &DB{}

type DB struct {
	pool *pgxpool.Pool
}

func New(url string) (*DB, error) {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, err
	}
	db := &DB{
		pool: pool,
	}

	return db, nil
}

func (db *DB) GetCommentsByArticleID(ctx context.Context, articleID int) ([]storage.Comment, error) {
	sql := `
SELECT
  id,
  article_id,
  parent_id,
  author,
  text,
  pub_time
FROM comments
WHERE article_id=$1
ORDER BY pub_time DESC;
`
	var comments []storage.Comment

	rows, err := db.pool.Query(ctx, sql, articleID)
	for rows.Next() {
		var c storage.Comment
		err = rows.Scan(
			&c.ID,
			&c.ArticleID,
			&c.ParentID,
			&c.Author,
			&c.Text,
			&c.PubTime,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (db *DB) AddComment(ctx context.Context, c storage.Comment) (int, error) {
	sql := `
INSERT INTO
  comments (
    article_id,
    parent_id,
    author,
    text,
    pub_time
  )
VALUES
  (
    $1, currval('comments_id_seq'), $2, $3, $4
  ) RETURNING id
`
	sqlForReply := `
INSERT INTO
  comments (
    article_id,
    parent_id,
    author,
    text,
    pub_time
  )
VALUES
  (
    $1, $2, $3, $4, $5
  ) RETURNING id
`
	var id int
	var err error
	if c.ParentID != 0 {
		err = db.pool.QueryRow(
			ctx,
			sqlForReply,
			c.ArticleID,
			c.ParentID,
			c.Author,
			c.Text,
			c.PubTime,
		).Scan(&id)
	} else {
		err = db.pool.QueryRow(
			ctx,
			sql,
			c.ArticleID,
			c.Author,
			c.Text,
			c.PubTime,
		).Scan(&id)
	}

	if err != nil {
		return 0, err
	}
	log.Println("Added comment", id)

	return id, nil
}

func (db *DB) DeleteComment(ctx context.Context, commentID int) error {
	sql := "DELETE FROM comments WHERE id = $1"

	r, err := db.pool.Exec(ctx, sql, commentID)
	if err != nil {
		return err
	}
	if r.RowsAffected() == 0 {
		log.Println("Comment ", commentID, "not found")
		return errors.New("not found")
	}
	log.Println("Deleted comment", commentID)

	return nil
}
