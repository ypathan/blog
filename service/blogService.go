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

func (s *BlogService) CreateBlog(blogReq model.BlogRequest) (*model.Blog, error) {
	blog := &model.Blog{
		Content: blogReq.Content,
		Title: blogReq.Title,
	}
	return s.repo.Create(blog)
}

func (s *BlogService) GetAllBlogs() ([]model.Blog, error) {
	return s.repo.FindAll()
}

func (s *BlogService) GetBlogByID(id int) (*model.Blog, error) {
	return s.repo.FindByID(id)
}

func (s *BlogService) UpdateBlog(id int, content string, title string) (*model.Blog, error) {
	return s.repo.Update(id, content, title)
}

func (s *BlogService) DeleteBlog(id int) error {
	return s.repo.Delete(id)
}
