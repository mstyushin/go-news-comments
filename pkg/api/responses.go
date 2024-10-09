package api

type CommentCreatedResponse struct {
	ID      int   `json:"id"`
	PubTime int64 `json:"pub_time"`
}
