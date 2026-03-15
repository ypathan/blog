package controller

import (
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"yousuf.xyz/blog/repository"
	"yousuf.xyz/blog/types"
)

type AuthController struct {
	repo *repository.AuthRepository
}

func NewAuthController(dbconn *sql.DB) *AuthController {
	return &AuthController{
		repo: repository.NewAuthRepository(dbconn),
	}
}

func (c *AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error reading payload", "error", err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	var user types.User
	if err := json.Unmarshal(body, &user); err != nil {
		slog.Error("error unmarshalling payload", "error", err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	//TODO: add user to DB logic
}

func (c *AuthController) LoginUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("error reading payload", "error", err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	var user types.User
	if err := json.Unmarshal(body, &user); err != nil {
		slog.Error("error unmarshalling payload", "error", err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	//TODO: login user logic
}
