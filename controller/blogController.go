package controller

import (
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"

	"yousuf.xyz/blog/repository"
	"yousuf.xyz/blog/types"
)

type BlogController struct {
	repo *repository.BlogRepository
}

func NewBlogController(dbconn *sql.DB) *BlogController {
	return &BlogController{
		repo: repository.NewBlogRepository(dbconn),
	}
}

func (s *BlogController) ServeIndex(w http.ResponseWriter, r *http.Request) {

	allBlogs, err := s.repo.FindAll()

	temp := template.Must(template.ParseFiles("static/index.gohtml", "static/ascii.html", "static/particles.html"))
	if err != nil {
		slog.Error("error getting all blogs", "error", err.Error())
	}

	ctx := map[string]any{
		"username": "ypathan",
		"allBlogs": allBlogs,
	}

	temp.Execute(w, ctx)
}

func (s *BlogController) ServeBlog(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("static/blog.gohtml", "static/ascii.html", "static/particles.html"))

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	blog, err := s.repo.FindByID(id)

	if err != nil {
		slog.Error("error getting all blogs", "error", err.Error())
	}

	ctx := map[string]any{
		"username": "ypathan",
		"blog":     blog,
	}

	temp.Execute(w, ctx)
}

func (c *BlogController) AddNewBlog(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read request body", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var blogReq types.BlogRequest
	if err := json.Unmarshal(body, &blogReq); err != nil {
		slog.Error("failed to unmarshal request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	blog, err := c.repo.Create(&types.Blog{
		Content: blogReq.Content,
		Title:   blogReq.Title,
	})

	if err != nil {
		slog.Error("failed to create blog", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		slog.Error("failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *BlogController) ViewAllBlogs(w http.ResponseWriter, r *http.Request) {
	blogs, err := c.repo.FindAll()
	if err != nil {
		slog.Error("failed to get all blogs", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blogs); err != nil {
		slog.Error("failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *BlogController) ViewBlog(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("invalid blog id", "id", idStr, "error", err)
		http.Error(w, "invalid blog id", http.StatusBadRequest)
		return
	}

	blog, err := c.repo.FindByID(id)
	if err != nil {
		slog.Error("failed to get blog", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if blog == nil {
		http.Error(w, "blog not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		slog.Error("failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *BlogController) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("invalid blog id", "id", idStr, "error", err)
		http.Error(w, "invalid blog id", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read request body", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var blogReq types.BlogRequest
	if err := json.Unmarshal(body, &blogReq); err != nil {
		slog.Error("failed to unmarshal request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	blog, err := c.repo.Update(id, blogReq.Content, blogReq.Title)
	if err != nil {
		slog.Error("failed to update blog", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		slog.Error("failed to encode response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *BlogController) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("invalid blog id", "id", idStr, "error", err)
		http.Error(w, "invalid blog id", http.StatusBadRequest)
		return
	}

	if err := c.repo.Delete(id); err != nil {
		slog.Error("failed to delete blog", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
