package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"text/template"

	"yousuf.xyz/blog/repository"
)

type AdminHandler struct {
	repo *repository.BlogRepository
}

func NewAdminHandler(_db *sql.DB) *AdminHandler {
	return &AdminHandler{
		repo: repository.NewBlogRepository(_db),
	}
}

func (a *AdminHandler) AdminPrivate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func (a *AdminHandler) AdminAddBlog(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("static/addblog.html", "static/particles.gohtml", "static/ascii.gohtml", "static/admintop.html"))
	ctx := map[string]any{}
	tmp.Execute(w, ctx)

}

func (a *AdminHandler) EditBlog(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("static/editblog.html", "static/particles.gohtml", "static/ascii.gohtml", "static/admintop.html"))
	ctx := map[string]any{}
	tmp.Execute(w, ctx)
}

func (a *AdminHandler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("static/adminDashboard.gohtml", "static/admintop.html", "static/particles.gohtml", "static/ascii.gohtml"))

	allBlogs, err := a.repo.FindAll()
	if err != nil {
		slog.Error("error fetching blogs for admin dashboard", "error")
	}

	ctx := map[string]any{
		"allBlogs": allBlogs,
	}
	tmp.Execute(w, ctx)
}
