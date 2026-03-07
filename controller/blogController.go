package controller

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"yousuf.xyz/blog/model"
	"yousuf.xyz/blog/service"
)

type BlogController struct {
	service *service.BlogService
}

func NewBlogController(service *service.BlogService) *BlogController {
	return &BlogController{service: service}
}

func (c *BlogController) AddNewBlog(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read request body", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var blogReq model.BlogRequest
	if err := json.Unmarshal(body, &blogReq); err != nil {
		slog.Error("failed to unmarshal request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	blog, err := c.service.CreateBlog(blogReq.Content)
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
	blogs, err := c.service.GetAllBlogs()
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

	blog, err := c.service.GetBlogByID(id)
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

	var blogReq model.BlogRequest
	if err := json.Unmarshal(body, &blogReq); err != nil {
		slog.Error("failed to unmarshal request body", "error", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	blog, err := c.service.UpdateBlog(id, blogReq.Content)
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

	if err := c.service.DeleteBlog(id); err != nil {
		slog.Error("failed to delete blog", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
