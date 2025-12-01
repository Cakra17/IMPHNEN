package store

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/Cakra17/imphnen/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3) RETURNING created_at 
	`
	err := r.db.QueryRowContext(ctx, query, user.ID, user.Email, user.PasswordHash).Scan(&user.Created_At)
	if err != nil {
		log.Printf("[ERROR] Failed to create user: %s", err.Error())
		return err
	}
	return nil
}

func (r *UserRepo) GetUserbyEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash FROM users WHERE email = $1`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err == sql.ErrNoRows {
		log.Printf("[ERROR] Failed to get user: %s", err.Error())
		return nil, errors.New("User not found")
	}
	if err != nil {
		log.Printf("[ERROR] Failed to get user: %s", err.Error())
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) Update(ctx context.Context, user models.User) error {
	query := `UPDATE TABLE users SET email = $1, password_hash = $2, WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, user.Email, user, user.PasswordHash, user.ID)
	if err != nil {
		log.Printf("[ERROR] Failed to update user: %s", err.Error())
		return err
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("[ERROR] Failed to delete user: %s", err.Error())
		return err
	}
	return nil
}
