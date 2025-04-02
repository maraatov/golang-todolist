package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	todo "todolists"
)

type AuthPostgres struct {
	db *pgxpool.Pool
}

func NewAuthPostgres(db *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	var id int
	query := "INSERT INTO users (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRow(context.Background(), query, user.Name, user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (todo.User, error) {
	var user todo.User
	query := "SELECT id, name, username FROM users WHERE username=$1 AND password_hash=$2"

	err := r.db.QueryRow(context.Background(), query, username, password).Scan(&user.Id, &user.Name, &user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}
