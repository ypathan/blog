package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"text/template"
	"time"

	"yousuf.xyz/blog/auth"
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

func (c *AuthController) ServeAdminLogin(w http.ResponseWriter, r *http.Request) {
	temp := template.Must(template.ParseFiles("static/adminLogin.gohtml", "static/ascii.gohtml", "static/particles.gohtml"))

	ctx := map[string]any{}

	temp.Execute(w, ctx)
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

	// for multiple users, can pass user object to check if user with same name exists
	exists := c.repo.CheckUserExists()
	if exists {
		http.Error(w, "user already exists", 500)
		return
	}

	user.Password = auth.HashPassword(user.Password)

	err = c.repo.RegisterUser(user)
	if err != nil {
		http.Error(w, "error registering user to db", 500)
		return
	}

	json.NewEncoder(w).Encode("register user success")

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

	// check username
	fetchedUser, err := c.repo.FindUserByUsername(user.Username)
	if err != nil {
		slog.Error("error fetching user for login", "error", err.Error())
		http.Error(w, "login failure, user does not exist", 500)
		return
	}

	// check password
	paswordMatched := auth.ComparePassword(user.Password, fetchedUser.Password)
	if !paswordMatched {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	session_token := auth.GenerateToken(32)
	csrf_token := auth.GenerateToken(32)

	//store session and csrf token in database
	fetchedUser.CSRFToken = csrf_token
	fetchedUser.SessionToken = session_token

	err = c.repo.UpdateUser(fetchedUser)
	if err != nil {
		slog.Error("error setting tokens for user", "error", err.Error())
		http.Error(w, "login failure, security error", 500)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session_token,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// should not expire, will always be there
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    fetchedUser.ID,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrf_token,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	json.NewEncoder(w).Encode("login success")
}

func (c *AuthController) LogoutUser(w http.ResponseWriter, r *http.Request) {

	// get user id
	user_id, err := r.Cookie("user_id")
	if err != nil {
		slog.Error("user id cookie not found", "error", err.Error())
		http.Error(w, "user id cookie not found", 404)
		return
	}

	err = c.repo.LogoutUser(user_id.Value)
	if err != nil {
		slog.Error("Error removing user tokens from DB", "error", err.Error())
		http.Error(w, "Internal Server Error During Logout", 404)
		return
	}

	// handling client side tokens
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})
	
	json.NewEncoder(w).Encode("logged out")
}
