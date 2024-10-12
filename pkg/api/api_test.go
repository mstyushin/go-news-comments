package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/mstyushin/go-news-comments/pkg/config"
	"github.com/mstyushin/go-news-comments/pkg/storage"
	"github.com/mstyushin/go-news-comments/pkg/storage/pg"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

const dbURL = "postgres://postgres@localhost:5433/comments?sslmode=disable"

var (
	s    *pg.DB
	pool *pgxpool.Pool
	api  *API
	ctx  context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	var err error
	s, err = pg.New(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	pool, err = pgxpool.Connect(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.DefaultConfig()
	api = New(cfg, s)
	api.initMux()

	os.Exit(m.Run())
}

func TestAPI_getComments(t *testing.T) {
	articleID := 1
	c := storage.Comment{
		ArticleID: articleID,
		Author:    "Bob",
		Text:      "Whatever",
		PubTime:   time.Now().Unix(),
	}

	cid, _ := api.db.AddComment(ctx, c)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/comments/by-articleid/%d", articleID), nil)
	rr := httptest.NewRecorder()

	api.mux.ServeHTTP(rr, req)
	assert.True(t, rr.Code == http.StatusOK)

	b, err := ioutil.ReadAll(rr.Body)
	assert.NoError(t, err, "should be able to read response body")

	var data []storage.Comment
	err = json.Unmarshal(b, &data)
	assert.NoError(t, err, "response should be deserializable")
	assert.NotEmpty(t, data, "should get non-empty slice of comments")
	assert.Equal(t, articleID, data[0].ArticleID, "should get comment for articleID=1")

	api.db.DeleteComment(ctx, cid)
}

func TestAPI_addComment(t *testing.T) {
	c := storage.Comment{
		ArticleID: 2,
		Author:    "Alice",
		Text:      "Hello world!",
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(c)

	req := httptest.NewRequest(http.MethodPost, "/comments", &buf)
	rr := httptest.NewRecorder()

	api.mux.ServeHTTP(rr, req)
	assert.True(t, rr.Code == http.StatusOK)

	b, err := ioutil.ReadAll(rr.Body)
	assert.NoError(t, err, "should be able to read response body")

	var data CommentCreatedResponse
	err = json.Unmarshal(b, &data)
	assert.NoError(t, err, "response should be deserializable")
	assert.NotEmpty(t, data, "should get non-nil response")
}
