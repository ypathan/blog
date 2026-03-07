package service

import (
	"database/sql"

	"yousuf.xyz/blog/model"
	"yousuf.xyz/blog/repository"
)

type BlogService struct {
	repo *repository.BlogRepository
}

func NewBlogService(db *sql.DB) *BlogService {
	return &BlogService{
		repo: repository.NewBlogRepository(db),
	}
}

func (s *BlogService) CreateBlog(content string) (*model.Blog, error) {
	blog := &model.Blog{
		Content: content,
	}
	return s.repo.Create(blog)
}

func (s *BlogService) GetAllBlogs() ([]model.Blog, error) {
	return s.repo.FindAll()
}

func (s *BlogService) GetBlogByID(id int) (*model.Blog, error) {
	return s.repo.FindByID(id)
}

func (s *BlogService) UpdateBlog(id int, content string) (*model.Blog, error) {
	return s.repo.Update(id, content)
}

func (s *BlogService) DeleteBlog(id int) error {
	return s.repo.Delete(id)
}
