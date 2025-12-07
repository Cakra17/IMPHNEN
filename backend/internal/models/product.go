package models

import "time"

type Product struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name" db:"name"`
	Price     float64   `json:"price" db:"price"`
	Stock     int       `json:"stock" db:"stock"`
	ImageURL  string    `json:"image_url" db:"image_url"`
	PublicID  string    `json:"public_id" db:"public_id"`
	CreatedAt time.Time `json:"created_at" db:"created-at"`
}

type ProductListResponse struct {
	Products []Product `json:"products"`
}
