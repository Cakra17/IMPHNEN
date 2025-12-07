package store

import (
	"context"
	"database/sql"
	"log"

	"github.com/Cakra17/imphnen/internal/models"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) ProductRepo {
	return ProductRepo{db: db}
}

func (r *ProductRepo) AddProduct(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO 
			products (id, user_id, name, price, stock, image_url, public_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		product.ID, product.UserID,
		product.Name, product.Price,
		product.Stock, product.ImageURL,
		product.PublicID,
	).Scan(&product.CreatedAt)
	if err != nil {
		log.Printf("[ERROR] Failed to add product: %s", err.Error())
		return err
	}
	return nil
}

func (r *ProductRepo) GetUserProductsPaginated(ctx context.Context, userID string, page, perPage uint) ([]models.Product, uint, error) {
	offset := (page - 1) * perPage
	query := `
		SELECT id, user_id, name, price, stock, image_url, public_id
		FROM products
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, perPage, offset)
	if err != nil {
		log.Printf("[ERROR] Failed to get products: %s", err.Error())
		return nil, 0, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID, &product.UserID,
			&product.Name, &product.Price,
			&product.Stock, &product.ImageURL,
			&product.PublicID,
		)
		if err != nil {
			log.Printf("[ERROR] Failed to scan receipt: %s", err.Error())
			return nil, 0, err
		}
		products = append(products, product)
	}

	// Get total count
	var totalCount uint
	countQuery := `SELECT COUNT(*) FROM products WHERE user_id = $1`
	err = r.db.QueryRowContext(ctx, countQuery, userID).Scan(&totalCount)
	if err != nil {
		log.Printf("[ERROR] Failed to get total count: %s", err.Error())
		return nil, 0, err
	}

	return products, totalCount, nil
}

func (r *ProductRepo) GetProductByID(ctx context.Context, productID string) (models.Product, error) {
	query := `
		SELECT id, user_id, name, price, stock, image_url, public_id
		FROM products
		WHERE id = $1
		LIMIT 1
	`
	var product models.Product
	err := r.db.QueryRowContext(ctx, query, productID).Scan(
		&product.ID, &product.UserID,
		&product.Name, &product.Price,
		&product.Stock, &product.ImageURL,
		&product.PublicID,
	)
	if err != nil {
		log.Printf("[ERROR] Failed to get product: %s", err.Error())
		return models.Product{}, err
	}
	return product, nil
}

func (r *ProductRepo) UpdateProduct(ctx context.Context, product models.Product) error {
	query := `
		UPDATE products 
		SET name=$1, price=$2, stock=$3, image_url=$4, public_id=$5
		WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query, product.Name, product.Price, product.Stock, product.ImageURL, product.PublicID, product.ID)
	if err != nil {
		log.Printf("[ERROR] Failed to update product: %s", err.Error())
		return err
	}
	return nil
}

func (r *ProductRepo) DeleteProduct(ctx context.Context, productID string) error {
	query := `
		DELETE FROM products WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, productID)
	if err != nil {
		log.Printf("[ERROR] Failed to update product: %s", err.Error())
		return err
	}
	return nil
}
