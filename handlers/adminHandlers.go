package handlers

import (
	"database/sql"
	"net/http"
	"text/template"
)

type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(_db *sql.DB) *AdminHandler{
	return &AdminHandler{
		db: _db,
	}
}

func (a *AdminHandler) AdminPrivate(w http.ResponseWriter , r *http.Request) {
	w.Write([]byte("hello world"))
}

func (a *AdminHandler) AdminAddBlog(w http.ResponseWriter, r *http.Request){
	tmp := template.Must(template.ParseFiles("static/addblog.html", "static/particles.gohtml", "static/ascii.gohtml"))
	ctx := map[string]any{}
	tmp.Execute(w, ctx)
	
}


func (a *AdminHandler) AdminDashboard(w http.ResponseWriter, r *http.Request){
	tmp := template.Must(template.ParseFiles("static/adminDashboard.html", "static/particles.gohtml", "static/ascii.gohtml"))
	ctx := map[string]any{}
	tmp.Execute(w, ctx)
	
}


