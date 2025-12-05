package models

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID         string      `json:"id" db:"id"`
	UserID     string      `json:"user_id" db:"user_id"`
	CustomerID string      `json:"customer_id" db:"customer_id"`
	TotalPrice float64     `json:"total_price" db:"total_price"`
	Status     OrderStatus `json:"status" db:"status"`
	OrderDate  time.Time   `json:"order_date" db:"order_date"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
	OrderItems []OrderItem `json:"order_items,omitempty" db:"-"`
	Customer   *Customer   `json:"customer,omitempty" db:"-"`
}

type OrderItem struct {
	ID         string    `json:"id" db:"id"`
	OrderID    string    `json:"order_id" db:"order_id"`
	ProductID  string    `json:"product_id" db:"product_id"`
	Quantity   int       `json:"quantity" db:"quantity"`
	TotalPrice float64   `json:"total_price" db:"total_price"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	Product    *Product  `json:"product,omitempty" db:"-"`
}

type Customer struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Address   string    `json:"address" db:"address"`
	Phone     string    `json:"phone" db:"phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateCustomerRequest struct {
	ID      int    `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
}

type CreateOrderRequest struct {
	CustomerID int                      `json:"customer_id" validate:"required"`
	Items      []CreateOrderItemRequest `json:"items" validate:"required,min=1,"`
}

type CreateTelegramOrderRequest struct {
	MerchantID string                   `json:"merchant_id"`
	CustomerID int                      `json:"customer_id" validate:"required"`
	Items      []CreateOrderItemRequest `json:"items" validate:"required,"`
}

type CreateOrderItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required,oneof=pending confirmed cancelled"`
}

type OrderFilter struct {
	UserID     string
	CustomerID *string
	Status     *OrderStatus
	Page       uint
	PerPage    uint
}

type OrderListResponse struct {
	Orders  []Order `json:"orders"`
	Total   uint    `json:"total"`
	Page    uint    `json:"page"`
	PerPage uint    `json:"per_page"`
}
