package handlers

import (
	"database/sql"
	"net/http"
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
