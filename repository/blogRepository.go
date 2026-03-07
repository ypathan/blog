package repository

import (
	"database/sql"
	"log/slog"

	"yousuf.xyz/blog/model"
)

type BlogRepository struct {
	db *sql.DB
}

func NewBlogRepository(db *sql.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) Create(blog *model.Blog) (*model.Blog, error) {
	result, err := r.db.Exec("INSERT INTO blog (content) VALUES (?)", blog.Content)
	if err != nil {
		slog.Error("failed to insert blog", "error", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("failed to get last insert id", "error", err)
		return nil, err
	}

	createdBlog, err := r.FindByID(int(id))
	if err != nil {
		slog.Error("failed to fetch created blog", "error", err)
		return nil, err
	}

	return createdBlog, nil
}

func (r *BlogRepository) FindByID(id int) (*model.Blog, error) {
	row := r.db.QueryRow("SELECT id, created_at, modified_at, is_deleted, content FROM blog WHERE id = ? AND is_deleted = 0", id)

	var blog model.Blog
	err := row.Scan(&blog.ID, &blog.CreatedAt, &blog.ModifiedAt, &blog.IsDeleted, &blog.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		slog.Error("failed to scan blog", "error", err)
		return nil, err
	}

	return &blog, nil
}

func (r *BlogRepository) FindAll() ([]model.Blog, error) {
	rows, err := r.db.Query("SELECT id, created_at, modified_at, is_deleted, content FROM blog WHERE is_deleted = 0")
	if err != nil {
		slog.Error("failed to query blogs", "error", err)
		return nil, err
	}
	defer rows.Close()

	var blogs []model.Blog
	for rows.Next() {
		var blog model.Blog
		err := rows.Scan(&blog.ID, &blog.CreatedAt, &blog.ModifiedAt, &blog.IsDeleted, &blog.Content)
		if err != nil {
			slog.Error("failed to scan blog row", "error", err)
			continue
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func (r *BlogRepository) Update(id int, content string) (*model.Blog, error) {
	_, err := r.db.Exec("UPDATE blog SET content = ? WHERE id = ?", content, id)
	if err != nil {
		slog.Error("failed to update blog", "error", err)
		return nil, err
	}

	return r.FindByID(id)
}

func (r *BlogRepository) Delete(id int) error {
	_, err := r.db.Exec("UPDATE blog SET is_deleted = 1 WHERE id = ?", id)
	if err != nil {
		slog.Error("failed to delete blog", "error", err)
		return err
	}
	return nil
}
