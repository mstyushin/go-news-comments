package api

import (
	"encoding/json"
	"go-news-comments/pkg/storage"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func (api *API) getComments(w http.ResponseWriter, r *http.Request) {
	// TODO implement pagination
	s := mux.Vars(r)["id"]
	articleID, _ := strconv.Atoi(s)
	comments, err := api.db.GetCommentsByArticleID(r.Context(), articleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (api *API) addComment(w http.ResponseWriter, r *http.Request) {
	var comment storage.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comment.PubTime = time.Now().Unix()
	cid, err := api.db.AddComment(r.Context(), comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CommentCreatedResponse{
		ID:      cid,
		PubTime: comment.PubTime,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}
