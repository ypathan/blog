package repository

import "database/sql"

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(_db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: _db,
	}
}

func (r *AuthRepository) RegisterUser() {

}
