package models

import "time"

type Transaction struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	Type            string    `json:"type" db:"type"`
	Source          string    `json:"source" db:"source"`
	Amount          float64   `json:"amount" db:"amount"`
	TransactionDate time.Time `json:"transaction_date" db:"transaction_date"`
	ReceiptID       *string   `json:"receipt_id" db:"receipt_id"`
	OrderID         *string   `json:"order_id" db:"order_id"`
	CreatedAt       time.Time `json:"created_at,omitempty" db:"created_at"`
}

type CreateTransactionPayload struct {
	Type            string  `json:"type" validate:"required,oneof=expense income"`
	Source          string  `json:"source" validate:"required,oneof=receipt bot manual"`
	Amount          float64 `json:"amount" validate:"required,gte=0"`
	TransactionDate string  `json:"transaction_date" validate:"required"`
	ReceiptID       *string `json:"receipt_id,omitempty"`
	OrderID         *string `json:"order_id,omitempty"`
}

type TransactionResponse struct {
	Transaction Transaction `json:"transaction"`
}

type TransactionListResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type TransactionStats struct {
	TotalIncome      float64 `json:"total_income"`
	TotalExpense     float64 `json:"total_expense"`
	NetAmount        float64 `json:"net_amount"`
	TransactionCount int64   `json:"transaction_count"`
	AverageAmount    float64 `json:"average_amount"`
}

type TransactionAnalytics struct {
	Stats        TransactionStats `json:"stats"`
	Transactions []Transaction    `json:"transactions,omitempty"`
}
