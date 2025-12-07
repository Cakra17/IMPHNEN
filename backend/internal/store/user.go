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
		INSERT INTO 
			users (id, email, password_hash, first_name, last_name, store_name) 
			VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at 
	`
	err := r.db.QueryRowContext(
		ctx, query,
		user.ID, user.Email,
		user.PasswordHash, user.FirstName,
		user.LastName, user.StoreName,
	).Scan(&user.Created_At)
	if err != nil {
		log.Printf("[ERROR] Failed to create user: %s", err.Error())
		return err
	}
	return nil
}

func (r *UserRepo) GetUserbyEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
	SELECT 
		id, email, password_hash, first_name, last_name, store_name 
	FROM users 
	WHERE email = $1`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.StoreName)
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

func (r *UserRepo) UpdateUser(ctx context.Context, id string, user models.User) error {
	query := `
		UPDATE users SET first_name = $1, last_name = $2, store_name = $3 WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, user.FirstName, user.LastName, user.StoreName, id)
	if err != nil {
		log.Printf("[ERROR] Failed to update user: %s", err.Error())
		return err
	}
	return nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("[ERROR] Failed to delete user: %s", err.Error())
		return err
	}
	return nil
}

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]models.Merchant, error) {
	query := `
		SELECT id, store_name 
		FROM users 
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("[ERROR] Failed to get all users: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var merchants []models.Merchant
	for rows.Next() {
		var merchant models.Merchant
		err := rows.Scan(&merchant.MerchantID, &merchant.MerchantName)
		if err != nil {
			log.Printf("[ERROR] Failed to scan user: %s", err.Error())
			return nil, err
		}
		merchants = append(merchants, merchant)
	}

	if err = rows.Err(); err != nil {
		log.Printf("[ERROR] Failed to iterate users: %s", err.Error())
		return nil, err
	}

	return merchants, nil
}
