package auth

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request){
	body, err := io.ReadAll(r.Body)

	if err != nil {
		slog.Error("error reading body", "error", err.Error())
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		slog.Error("error unmarshalling the json payload", "error", err.Error())
	}



	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)

}

func LoginUser(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("hello world from login"))	
}

