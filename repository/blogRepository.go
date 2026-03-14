package repository

import (
	"database/sql"
	"log/slog"
	"time"

	"yousuf.xyz/blog/model"
)

type BlogRepository struct {
	db *sql.DB
}

func NewBlogRepository(db *sql.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) Create(blog *model.Blog) (*model.Blog, error) {
	var id int
	err := r.db.QueryRow("INSERT INTO blog (content, title) VALUES ($1, $2) RETURNING id", blog.Content, blog.Title).Scan(&id)
	if err != nil {
		slog.Error("failed to insert blog", "error", err)
		return nil, err
	}

	createdBlog, err := r.FindByID(id)
	if err != nil {
		slog.Error("failed to fetch created blog", "error", err)
		return nil, err
	}

	return createdBlog, nil
}

func (r *BlogRepository) FindByID(id int) (*model.Blog, error) {
	row := r.db.QueryRow("SELECT id, created_at, modified_at, is_deleted, content, title FROM blog WHERE id = $1 AND is_deleted = false", id)

	var blog model.Blog
	err := row.Scan(&blog.ID, &blog.CreatedAt, &blog.ModifiedAt, &blog.IsDeleted, &blog.Content, &blog.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("failed to scan blog", "error", err)
		return nil, err
	}

	t, _ := time.Parse(time.RFC3339Nano, blog.CreatedAt)
	blog.CreatedAt = t.Format("2006-01-02")

	return &blog, nil
}

func (r *BlogRepository) FindAll() ([]model.Blog, error) {
	rows, err := r.db.Query("SELECT id, created_at, modified_at, is_deleted, content, title FROM blog WHERE is_deleted = false")
	if err != nil {
		slog.Error("failed to query blogs", "error", err)
		return nil, err
	}
	defer rows.Close()
	var blogs []model.Blog
	for rows.Next() {
		var blog model.Blog
		err := rows.Scan(&blog.ID, &blog.CreatedAt, &blog.ModifiedAt, &blog.IsDeleted, &blog.Content, &blog.Title)
		if err != nil {
			slog.Error("failed to scan blog row", "error", err)
			continue
		}

		// format content
		if len(blog.Content) > 500 {
			blog.Content = blog.Content[:500]
		}

		// format date
		t, _ := time.Parse(time.RFC3339Nano, blog.CreatedAt)
		blog.CreatedAt = t.Format("2006-01-02")

		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func (r *BlogRepository) Update(id int, content string, title string) (*model.Blog, error) {

	_, err := r.db.Exec("UPDATE blog SET content = $1, title = $2  WHERE id = $3", content, title, id)
	if err != nil {
		slog.Error("failed to update blog", "error", err)
		return nil, err
	}

	return r.FindByID(id)
}

func (r *BlogRepository) Delete(id int) error {
	_, err := r.db.Exec("UPDATE blog SET is_deleted = true WHERE id = $1", id)
	if err != nil {
		slog.Error("failed to delete blog", "error", err)
		return err
	}
	return nil
}
