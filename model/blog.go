package model

type Blog struct {
	ID         int    `json:"id"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
	IsDeleted  bool    `json:"is_deleted"`
	Content    string `json:"content"`
}

type BlogRequest struct {
	Content string `json:"content"`
}
