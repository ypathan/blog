package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Blog struct {
	ID         int    `json:"id"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
	IsDeleted  int    `json:"is_deleted"`
	Content    string `json:"content"`
}

type BlogContent struct {
	Content string `json:"content"`
}


type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) ViewAllBlogs(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query("select * from blog where is_deleted = 0")
	if err != nil {
		log.Println("error: ", err.Error())
	}
	defer rows.Close()

	var allblogs []Blog
	for rows.Next() {
		var _blog Blog

		err := rows.Scan(&_blog.ID, &_blog.CreatedAt, &_blog.ModifiedAt, &_blog.IsDeleted, &_blog.Content)
		if err != nil {
			log.Println("error: ", err.Error())
		}

		allblogs = append(allblogs, _blog)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allblogs); err != nil {
		log.Println("error:", err.Error())
	}
}

func (s *Service) AddNewBlog(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error: ", err.Error())
	}

	var blogContent BlogContent
	json.Unmarshal(body, &blogContent)

	fmt.Println(blogContent.Content)

	res, err := s.db.Exec("insert into blog (content) values (?)", blogContent.Content)
	if err != nil {
		log.Println("error: ", err.Error())
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("error: ", err.Error())
	}

	row, err := s.db.Query("select * from blog where id = ? and is_deleted = 0", id)
	if err != nil {
		log.Println("error: ", err.Error())
	}

	var blog Blog
	for row.Next() {
		row.Scan(&blog.ID, &blog.CreatedAt, &blog.ModifiedAt, &blog.IsDeleted, &blog.Content)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)

}

func (s *Service) UpdateBlog(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	// check if such record exists
	rows, err := s.db.Query("select * from blog where id = ? and is_deleted = 0", id)
	if err != nil {
		log.Println("error: ", err.Error())
	}

	var blog Blog
	for rows.Next() {
		rows.Scan(&blog.ID, &blog.CreatedAt, &blog.ModifiedAt, &blog.IsDeleted, &blog.Content)
	}

	// Unmarshalling the body for update
	var blogcontent BlogContent
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error: ", err.Error())
	}

	json.Unmarshal(body, &blogcontent)

	s.db.Exec("update blog set content = ? where id = ?", blogcontent.Content, id)

	updated_rows, err := s.db.Query("select * from blog where id = ? and is_deleted = 0", id)
	if err != nil {
		log.Println("error: ", err.Error())
	}
	for updated_rows.Next() {
		updated_rows.Scan(&blog.ID, &blog.CreatedAt, &blog.ModifiedAt, &blog.IsDeleted, &blog.Content)
	}

	if err := json.NewEncoder(w).Encode(blog); err != nil {
		log.Println("error: ", err.Error())
	}
}

func (s *Service) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	_, err := s.db.Exec("update blog set is_deleted = 1 where id = ?", id)
	if err != nil {
		log.Println("error: ", err.Error())
	}

	row, err := s.db.Query("select * from blog where id  = ? and is_deleted = 0", id)
	if err != nil {
		log.Println("error: ", err.Error())
	}

	var blog Blog
	for row.Next() {
		err := row.Scan(&blog.ID, &blog.CreatedAt, &blog.ModifiedAt, &blog.IsDeleted, &blog.Content)
		if err != nil {
			log.Println("error: ", err.Error())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func (s *Service) View(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	rows, err := s.db.Query("select * from blog where id = ? and is_deleted = 0", id)
	if err != nil {
		log.Println("error: ", err.Error())
	}
	defer rows.Close()

	var _blog Blog
	for rows.Next() {
		err := rows.Scan(&_blog.ID, &_blog.CreatedAt, &_blog.ModifiedAt, &_blog.IsDeleted, &_blog.Content)
		if err != nil {
			log.Println("error: ", err.Error())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(_blog); err != nil {
		log.Println("error:", err.Error())
	}
}
