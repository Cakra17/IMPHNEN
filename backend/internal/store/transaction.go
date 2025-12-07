package store

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Cakra17/imphnen/internal/models"
)

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) TransactionRepo {
	return TransactionRepo{db: db}
}

func (s *TransactionRepo) AddTransaction(ctx context.Context, transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions (id, user_id, type, source, amount, transaction_date, receipt_id, order_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at
	`

	err := s.db.QueryRowContext(
		ctx, query,
		transaction.ID, transaction.UserID,
		transaction.Type, transaction.Source,
		transaction.Amount, transaction.TransactionDate,
		transaction.ReceiptID, transaction.OrderID,
	).Scan(&transaction.CreatedAt)

	if err != nil {
		log.Printf("[ERROR] Failed to create transaction: %s", err.Error())
		return err
	}

	return nil
}

func (s *TransactionRepo) GetTransactionsByDate(ctx context.Context, userID string, date time.Time) ([]models.Transaction, error) {
	query := `
		SELECT id, user_id, type, source, amount, transaction_date, receipt_id, order_id, created_at
		FROM transactions
		WHERE user_id = $1 AND DATE(transaction_date) = $2
		ORDER BY transaction_date DESC
	`
	rows, err := s.db.QueryContext(ctx, query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Type,
			&transaction.Source,
			&transaction.Amount,
			&transaction.TransactionDate,
			&transaction.ReceiptID,
			&transaction.OrderID,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (s *TransactionRepo) GetTransactionsByRange(ctx context.Context, userID string, startDate, endDate time.Time) ([]models.Transaction, error) {
	query := `
		SELECT id, user_id, type, source, amount, transaction_date, receipt_id, order_id, created_at
		FROM transactions
		WHERE user_id = $1 AND DATE(transaction_date) BETWEEN $2 AND $3
		ORDER BY transaction_date DESC
	`
	rows, err := s.db.QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Type,
			&transaction.Source,
			&transaction.Amount,
			&transaction.TransactionDate,
			&transaction.ReceiptID,
			&transaction.OrderID,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (s *TransactionRepo) GetTransactionsByDays(ctx context.Context, userID string, days int) ([]models.Transaction, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	return s.GetTransactionsByRange(ctx, userID, startDate, endDate)
}

func (s *TransactionRepo) GetTransactionStats(ctx context.Context, userID string, startDate, endDate time.Time) (models.TransactionStats, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as total_expense,
			COUNT(*) as transaction_count,
			COALESCE(AVG(amount), 0) as average_amount
		FROM transactions
		WHERE user_id = $1 AND DATE(transaction_date) BETWEEN $2 AND $3
	`

	var stats models.TransactionStats
	err := s.db.QueryRowContext(ctx, query, userID, startDate, endDate).Scan(
		&stats.TotalIncome,
		&stats.TotalExpense,
		&stats.TransactionCount,
		&stats.AverageAmount,
	)
	if err != nil {
		return models.TransactionStats{}, err
	}

	stats.NetAmount = stats.TotalIncome - stats.TotalExpense
	return stats, nil
}

func (s *TransactionRepo) GetTransactionStatsByDays(ctx context.Context, userID string, days int) (models.TransactionStats, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	return s.GetTransactionStats(ctx, userID, startDate, endDate)
}

func (s *TransactionRepo) GetTransactionsByType(
	ctx context.Context, userID string, transactionType string, startDate, endDate time.Time,
) ([]models.Transaction, error) {
	query := `
		SELECT id, user_id, type, source, amount, transaction_date, receipt_id, created_at
		FROM transactions
		WHERE user_id = $1 AND type = $2 AND DATE(transaction_date) BETWEEN $3 AND $4
		ORDER BY transaction_date DESC
	`
	rows, err := s.db.QueryContext(ctx, query, userID, transactionType, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Type,
			&transaction.Source,
			&transaction.Amount,
			&transaction.TransactionDate,
			&transaction.ReceiptID,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func (s *TransactionRepo) GetTransactionsBySource(
	ctx context.Context, userID string, source string, startDate, endDate time.Time,
) ([]models.Transaction, error) {
	query := `
		SELECT id, user_id, type, source, amount, transaction_date, receipt_id, order_id, created_at
		FROM transactions
		WHERE user_id = $1 AND source = $2 AND DATE(transaction_date) BETWEEN $3 AND $4
		ORDER BY transaction_date DESC
	`
	rows, err := s.db.QueryContext(ctx, query, userID, source, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.Type,
			&transaction.Source,
			&transaction.Amount,
			&transaction.TransactionDate,
			&transaction.ReceiptID,
			&transaction.OrderID,
			&transaction.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
