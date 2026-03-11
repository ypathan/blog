package model

type Blog struct {
	ID         int    `json:"id"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
	IsDeleted  bool    `json:"is_deleted"`
	Content    string `json:"content"`
	Title 	string 	`json:"title"`
}

type BlogRequest struct {
	Title 	string 	`json:"title"`
	Content string `json:"content"`
}
