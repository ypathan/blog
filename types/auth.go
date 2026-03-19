package types

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInternal struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	CSRFToken    string `json:"csrf_token"`
	SessionToken string `json:"session_token"`
}
