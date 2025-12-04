package models

import "time"

type Receipt struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	StoreName  string    `json:"store_name" db:"store_name"`
	TotalItems uint32    `json:"total_items" db:"total_items"`
	TotalPrice float64   `json:"total" db:"total_price"`
	ImageURL   string    `json:"image_url" db:"image_url"`
	CreatedAt  time.Time `json:"created_at,omitempty" db:"created_at"`
}

type ReceiptItem struct {
	ID        string    `json:"id" db:"id"`
	ReceiptID string    `json:"receipt_id" db:"receipt_id"`
	Name      string    `json:"name" db:"name"`
	Price     float64   `json:"price" db:"price"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

type IncomeAndExpenses struct {
	Incomes  float64 `json:"incomes"`
	Expenses float64 `json:"expenses"`
}

type CreateReceiptPayload struct {
	StoreName  string                     `json:"store_name" validate:"required"`
	TotalItems uint32                     `json:"total_items" validate:"required"`
	TotalPrice float64                    `json:"total" validate:"required"`
	ImageURL   string                     `json:"image_url" validate:"required"`
	Items      []CreateReceiptItemPayload `json:"items" validate:"required,min=1"`
}

type CreateReceiptItemPayload struct {
	Name  string  `json:"name" validate:"required"`
	Price float64 `json:"price" validate:"required"`
}

type ReceiptResponse struct {
	Receipt Receipt       `json:"receipt"`
	Items   []ReceiptItem `json:"items,omitempty"`
}

type ReceiptListResponse struct {
	Receipts []Receipt `json:"receipts"`
}

type ItemsResponse struct {
	Items []ReceiptItem `json:"items"`
}

type InvoiceItems struct {
	Description string  `json:"description"`
	Total       float64 `json:"total"`
}

type Issuer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type OCR struct {
	InvoiceDate  string         `json:"invoice_date"`
	Subtotal     float64        `json:"subtotal"`
	Tax          float64        `json:"tax"`
	Total        float64        `json:"total"`
	Issuer       Issuer         `json:"issuer"`
	InvoiceItems []InvoiceItems `json:"invoice_items"`
}
