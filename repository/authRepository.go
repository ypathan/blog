package repository

import (
	"database/sql"
	"log/slog"
	"strconv"

	"yousuf.xyz/blog/types"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(_db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: _db,
	}
}

func (r *AuthRepository) CheckUserExists() bool {
	// i am the only user for now
	var exists bool
	err := r.db.QueryRow("select exists ( select username from users where username = 'ypathan' ) ").Scan(&exists)
	if err != nil {
		slog.Error("error checking user exists", "error", err.Error())
	}
	return exists
}

func (r *AuthRepository) RegisterUser(user types.User) error {
	var id int
	err := r.db.QueryRow("insert into users (username, password) values ($1, $2) returning id", user.Username, user.Password).Scan(&id)
	return err
}

func (r *AuthRepository) FindUserByUsername(username string) (types.UserInternal, error) {
	row := r.db.QueryRow("select id, username, password, session_token, csrf_token from users where username = $1 ", username)

	var user types.UserInternal
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.SessionToken, &user.CSRFToken)

	return user, err
}

func (r *AuthRepository) UpdateUser(user types.UserInternal) error {
	_, err := r.db.Exec("update users set session_token = $1 , csrf_token = $2 where username = $3", user.SessionToken, user.CSRFToken, user.Username)
	return err
}

func (r *AuthRepository) LogoutUser(user_id string) error {

	id, err := strconv.Atoi(user_id)
	if err != nil {

	}

	_, err = r.db.Exec("update users set session_token = '' , csrf_token = '' where id = $1", id)
	return err
}
