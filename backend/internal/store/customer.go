package store

import (
	"context"
	"database/sql"
	"log"

	"github.com/Cakra17/imphnen/internal/models"
)

type CustomerRepo struct {
	db *sql.DB
}

func NewCustomerRepo(db *sql.DB) CustomerRepo {
	return CustomerRepo{db: db}
}

func (r *CustomerRepo) CreateCustomer(ctx context.Context, customer *models.Customer) error {
	query := `
		INSERT INTO 
			customers (id, name, address, phone)
			VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		customer.ID, customer.Name,
		customer.Address, customer.Phone,
	).Scan(&customer.CreatedAt)
	if err != nil {
		log.Printf("[ERROR] Failed to create customer: %s", err.Error())
		return err
	}
	return nil
}

func (r *CustomerRepo) GetCustomerByID(ctx context.Context, customerID int) (models.Customer, error) {
	query := `
		SELECT id, name, address, phone, created_at
		FROM customers
		WHERE id = $1
		LIMIT 1
	`
	var customer models.Customer
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(
		&customer.ID, &customer.Name,
		&customer.Address, &customer.Phone,
		&customer.CreatedAt,
	)
	if err != nil {
		log.Printf("[ERROR] Failed to get customer: %s", err.Error())
		return models.Customer{}, err
	}
	return customer, nil
}

func (r *CustomerRepo) UpdateCustomer(ctx context.Context, customer models.Customer) error {
	query := `
		UPDATE customers 
		SET name=$1, address=$2, phone=$3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, customer.Name, customer.Address, customer.Phone, customer.ID)
	if err != nil {
		log.Printf("[ERROR] Failed to update customer: %s", err.Error())
		return err
	}
	return nil
}

func (r *CustomerRepo) DeleteCustomer(ctx context.Context, customerID int) error {
	query := `
		DELETE FROM customers WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, customerID)
	if err != nil {
		log.Printf("[ERROR] Failed to delete customer: %s", err.Error())
		return err
	}
	return nil
}
