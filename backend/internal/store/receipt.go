package store

import (
	"context"
	"database/sql"
	"log"

	"github.com/Cakra17/imphnen/internal/models"
)

type ReceiptRepo struct {
	db *sql.DB
}

func NewReceiptRepo(db *sql.DB) ReceiptRepo {
	return ReceiptRepo{db: db}
}

func (r *ReceiptRepo) Create(ctx context.Context, receipt models.Receipt) error {
	query := `
		INSERT INTO 
			receipts (id, user_id, total_items, total_price, store_name, image_url) 
			VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at
	`
	err := r.db.QueryRowContext(ctx, query,
		receipt.ID, receipt.UserID, receipt.TotalItems,
		receipt.TotalPrice, receipt.StoreName,
		receipt.ImageURL,
	).Scan(&receipt.CreatedAt)
	if err != nil {
		log.Printf("[ERROR] Failed to create receipt: %s", err.Error())
		return err
	}
	return nil
}

func (r *ReceiptRepo) CreateItems(ctx context.Context, items []models.ReceiptItem) error {
	query := `
		INSERT INTO receipt_items (id, receipt_id, name, price) 
		VALUES ($1, $2, $3, $4)
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to begin transaction: %s", err.Error())
		return err
	}
	defer tx.Rollback()

	for _, item := range items {
		_, err := tx.ExecContext(ctx, query, item.ID, item.ReceiptID, item.Name, item.Price)
		if err != nil {
			log.Printf("[ERROR] Failed to create receipt item: %s", err.Error())
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[ERROR] Failed to commit transaction: %s", err.Error())
		return err
	}

	return nil
}

func (r *ReceiptRepo) GetReceiptsPaginate(ctx context.Context, userID string, page, perPage uint) ([]models.Receipt, uint, error) {
	offset := (page - 1) * perPage

	query := `
		SELECT id, user_id, total_items, total_price, store_name, image_url, created_at
		FROM receipts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, perPage, offset)
	if err != nil {
		log.Printf("[ERROR] Failed to get receipts: %s", err.Error())
		return nil, 0, err
	}
	defer rows.Close()

	receipts := []models.Receipt{}
	for rows.Next() {
		var receipt models.Receipt
		err := rows.Scan(
			&receipt.ID, &receipt.UserID, &receipt.TotalItems,
			&receipt.TotalPrice, &receipt.StoreName,
			&receipt.ImageURL, &receipt.CreatedAt,
		)
		if err != nil {
			log.Printf("[ERROR] Failed to scan receipt: %s", err.Error())
			return nil, 0, err
		}
		receipts = append(receipts, receipt)
	}

	// Get total count
	var totalCount uint
	countQuery := `SELECT COUNT(*) FROM receipts WHERE user_id = $1`
	err = r.db.QueryRowContext(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		log.Printf("[ERROR] Failed to get total count: %s", err.Error())
		return nil, 0, err
	}

	return receipts, totalCount, nil
}

func (r *ReceiptRepo) GetReceiptByID(ctx context.Context, receiptID string, userID string) (*models.Receipt, error) {
	query := `
		SELECT id, user_id, total_items, total_price, store_name, image_url, created_at
		FROM receipts
		WHERE id = $1 AND user_id = $2
	`

	var receipt models.Receipt
	err := r.db.QueryRowContext(ctx, query, receiptID, userID).Scan(
		&receipt.ID, &receipt.UserID, &receipt.TotalItems,
		&receipt.TotalPrice, &receipt.StoreName,
		&receipt.ImageURL, &receipt.CreatedAt,
	)
	if err == sql.ErrNoRows {
		log.Printf("[ERROR] Receipt not found: %s", err.Error())
		return nil, err
	}
	if err != nil {
		log.Printf("[ERROR] Failed to get receipt: %s", err.Error())
		return nil, err
	}

	return &receipt, nil
}

func (r *ReceiptRepo) GetReceiptItemsByReceiptID(ctx context.Context, receiptID string) ([]models.ReceiptItem, error) {
	query := `
		SELECT id, receipt_id, name, price, created_at
		FROM receipt_items
		WHERE receipt_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, receiptID)
	if err != nil {
		log.Printf("[ERROR] Failed to get receipt items: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	items := []models.ReceiptItem{}
	for rows.Next() {
		var item models.ReceiptItem
		err := rows.Scan(&item.ID, &item.ReceiptID, &item.Name, &item.Price, &item.CreatedAt)
		if err != nil {
			log.Printf("[ERROR] Failed to scan receipt item: %s", err.Error())
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
