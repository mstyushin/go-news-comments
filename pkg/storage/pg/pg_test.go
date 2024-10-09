package pg

import (
	"context"
	"errors"
	"fmt"
	"go-news-comments/pkg/storage"
	"log"
	"os"
	"testing"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

const dbURL = "postgres://postgres@localhost:5432/comments?sslmode=disable"

var (
	s             *DB
	pool          *pgxpool.Pool
	ctx           context.Context
	testComment   storage.Comment
	commentAuthor = "Alice"
	commentText   = "something very thoughtful and interesting"
	pubTime       = int64(time.Now().Unix())
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	var err error
	s, err = New(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	pool, err = pgxpool.Connect(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	testComment = storage.Comment{
		ArticleID: 1,
		Author:    commentAuthor,
		Text:      commentText,
		PubTime:   pubTime,
	}

	os.Exit(m.Run())
}

func getCommentByID(id int) (storage.Comment, error) {
	sql := `
SELECT
  id,
  article_id,
  parent_id,
  author,
  text,
  pub_time
FROM comments
WHERE id=$1;
`
	var res []storage.Comment
	if err := pgxscan.Select(ctx, pool, &res, sql, id); err != nil {
		return storage.Comment{}, err
	}

	if len(res) == 1 {
		return res[0], nil
	}

	return storage.Comment{}, errors.New("unexpected search result")
}

func TestPG_AddComment(t *testing.T) {
	cid, err := s.AddComment(ctx, testComment)
	assert.NoError(t, err, "adding comment")

	_c, err := getCommentByID(cid)
	assert.NoError(t, err, fmt.Sprintf("searching for comment id=%d", cid))

	assert.Equal(t, commentAuthor, _c.Author)
	assert.Equal(t, commentText, _c.Text)
	assert.Equal(t, pubTime, _c.PubTime)

	s.DeleteComment(ctx, cid)
}

func TestPG_DeleteComment(t *testing.T) {
	cid, _ := s.AddComment(ctx, testComment)

	err := s.DeleteComment(ctx, cid)
	assert.NoError(t, err, "deleting comment")

	_, err = getCommentByID(cid)
	assert.EqualError(t, err, "unexpected search result")
}

func TestPG_GetComments(t *testing.T) {
	aid := 1
	latestTS := pubTime + 3
	comments := []storage.Comment{
		storage.Comment{
			ArticleID: aid,
			Author:    "Alice",
			Text:      "one",
			PubTime:   pubTime + 1,
		},
		storage.Comment{
			ArticleID: aid,
			Author:    "Bob",
			Text:      "two",
			PubTime:   latestTS,
		},
		storage.Comment{
			ArticleID: aid + 1,
			Author:    "Jane",
			Text:      "three",
			PubTime:   pubTime + 2,
		},
	}

	for _, c := range comments {
		s.AddComment(ctx, c)
	}

	_comments, err := s.GetCommentsByArticleID(ctx, aid)
	assert.NoError(t, err, "getting comments by articleID")
	assert.NotEmpty(t, _comments, "comments slice should not be empty")
	assert.Equal(t, 2, len(_comments), "should get slice of two comments")
	assert.Equal(t, latestTS, _comments[0].PubTime, "first comment in the slice should be the latest one")
}
