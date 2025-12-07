package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Cakra17/imphnen/internal/models"
	"github.com/lib/pq"
)

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) OrderRepo {
	return OrderRepo{db: db}
}

func (r *OrderRepo) CreateOrder(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to begin transaction: %s", err.Error())
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("[ERROR] Failed to rollback transaction: %s", rbErr.Error())
			}
		}
	}()

	productIDs := make([]string, len(items))
	requestedQuantities := make(map[string]int)
	for i, item := range items {
		productIDs[i] = item.ProductID
		requestedQuantities[item.ProductID] = item.Quantity
	}

	query := `
		SELECT id, stock, price 
		FROM products 
		WHERE id = ANY($1) AND user_id = $2
		FOR UPDATE
	`
	rows, err := tx.QueryContext(ctx, query, pq.Array(productIDs), order.UserID)
	if err != nil {
		log.Printf("[ERROR] Failed to lock products: %s", err.Error())
		return err
	}
	defer rows.Close()

	productStockMap := make(map[string]int)
	productPriceMap := make(map[string]float64)
	scannedProducts := 0

	for rows.Next() {
		var productID string
		var stock int
		var price float64
		if err := rows.Scan(&productID, &stock, &price); err != nil {
			log.Printf("[ERROR] Failed to scan product: %s", err.Error())
			return err
		}
		productStockMap[productID] = stock
		productPriceMap[productID] = price
		scannedProducts++
	}

	if scannedProducts != len(productIDs) {
		return fmt.Errorf("one or more products not found or don't belong to user")
	}

	for productID, requestedQty := range requestedQuantities {
		availableStock, exists := productStockMap[productID]
		if !exists {
			return fmt.Errorf("product %s not found", productID)
		}
		if availableStock < requestedQty {
			return fmt.Errorf("insufficient stock for product %s: requested %d, available %d",
				productID, requestedQty, availableStock)
		}
	}

	var customerExists bool
	err = tx.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM customers WHERE id = $1)`,
		order.CustomerID).Scan(&customerExists)
	if err != nil {
		log.Printf("[ERROR] Failed to check customer: %s", err.Error())
		return err
	}
	if !customerExists {
		return fmt.Errorf("customer not found")
	}

	var totalPrice float64
	for i := range items {
		price := productPriceMap[items[i].ProductID]
		items[i].TotalPrice = price * float64(items[i].Quantity)
		totalPrice += items[i].TotalPrice

		updateStockQuery := `
			UPDATE products 
			SET stock = stock - $1 
			WHERE id = $2 AND stock >= $1
		`
		result, err := tx.ExecContext(ctx, updateStockQuery, items[i].Quantity, items[i].ProductID)
		if err != nil {
			log.Printf("[ERROR] Failed to update product stock: %s", err.Error())
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("[ERROR] Failed to check rows affected: %s", err.Error())
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("insufficient stock for product %s (concurrent update)", items[i].ProductID)
		}
	}

	order.TotalPrice = totalPrice

	orderQuery := `
		INSERT INTO orders(id, user_id, customer_id, total_price, status, order_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`
	err = tx.QueryRowContext(
		ctx, orderQuery,
		order.ID, order.UserID, order.CustomerID,
		order.TotalPrice, order.Status,
		order.OrderDate, order.CreatedAt,
	).Scan(&order.CreatedAt)
	if err != nil {
		log.Printf("[ERROR] Failed to create order: %s", err.Error())
		return err
	}

	for _, item := range items {
		itemQuery := `
			INSERT INTO order_items(id, order_id, product_id, quantity, total_price, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err = tx.ExecContext(
			ctx, itemQuery,
			item.ID, order.ID, item.ProductID,
			item.Quantity, item.TotalPrice, item.CreatedAt,
		)
		if err != nil {
			log.Printf("[ERROR] Failed to create order item: %s", err.Error())
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[ERROR] Failed to commit transaction: %s", err.Error())
		return err
	}

	order.OrderItems = items
	return nil
}

func (r *OrderRepo) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	query := `
		SELECT 
			o.id, o.user_id, o.customer_id, o.total_price, o.status, o.order_date, o.created_at,
			c.id, c.name, c.address, c.phone, c.created_at
		FROM orders o
		LEFT JOIN customers c ON o.customer_id = c.id
		WHERE o.id = $1
	`

	var order models.Order
	var customer models.Customer
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&order.ID, &order.UserID, &order.CustomerID, &order.TotalPrice,
		&order.Status, &order.OrderDate, &order.CreatedAt,
		&customer.ID, &customer.Name, &customer.Address, &customer.Phone, &customer.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		log.Printf("[ERROR] Failed to get order: %s", err.Error())
		return nil, err
	}

	order.Customer = &customer

	items, err := r.getOrderItems(ctx, orderID)
	if err != nil {
		return nil, err
	}
	order.OrderItems = items

	return &order, nil
}

func (r *OrderRepo) getOrderItems(ctx context.Context, orderID string) ([]models.OrderItem, error) {
	query := `
		SELECT 
			oi.id, oi.order_id, oi.product_id, oi.quantity, oi.total_price, oi.created_at,
			p.id, p.name, p.price, p.stock, p.image_url
		FROM order_items oi
		LEFT JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = $1
		ORDER BY oi.created_at
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		log.Printf("[ERROR] Failed to get order items: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	items := []models.OrderItem{}
	for rows.Next() {
		var item models.OrderItem
		var product models.Product
		var productID, productName, productImageURL sql.NullString
		var productPrice sql.NullFloat64
		var productStock sql.NullInt32

		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.TotalPrice, &item.CreatedAt,
			&productID, &productName, &productPrice, &productStock, &productImageURL,
		)
		if err != nil {
			log.Printf("[ERROR] Failed to scan order item: %s", err.Error())
			return nil, err
		}

		if productID.Valid {
			product.ID = productID.String
			product.Name = productName.String
			product.Price = productPrice.Float64
			product.Stock = int(productStock.Int32)
			product.ImageURL = productImageURL.String
			item.Product = &product
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *OrderRepo) GetOrders(ctx context.Context, filter models.OrderFilter) ([]models.Order, uint, error) {
	whereConditions := []string{"o.user_id = $1"}
	args := []interface{}{filter.UserID}
	argCount := 1

	if filter.CustomerID != nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("o.customer_id = $%d", argCount))
		args = append(args, *filter.CustomerID)
	}

	if filter.Status != nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("o.status = $%d", argCount))
		args = append(args, *filter.Status)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	var totalCount uint
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM orders o WHERE %s`, whereClause)
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		log.Printf("[ERROR] Failed to get total count: %s", err.Error())
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PerPage
	argCount++
	limitArg := argCount
	argCount++
	offsetArg := argCount

	query := fmt.Sprintf(`
		SELECT 
			o.id, o.user_id, o.customer_id, o.total_price, o.status, o.order_date, o.created_at,
			c.id, c.name, c.address, c.phone, c.created_at
		FROM orders o
		LEFT JOIN customers c ON o.customer_id = c.id
		WHERE %s
		ORDER BY o.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, limitArg, offsetArg)

	args = append(args, filter.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("[ERROR] Failed to get orders: %s", err.Error())
		return nil, 0, err
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var order models.Order
		var customer models.Customer
		err := rows.Scan(
			&order.ID, &order.UserID, &order.CustomerID, &order.TotalPrice,
			&order.Status, &order.OrderDate, &order.CreatedAt,
			&customer.ID, &customer.Name, &customer.Address, &customer.Phone, &customer.CreatedAt,
		)
		if err != nil {
			log.Printf("[ERROR] Failed to scan order: %s", err.Error())
			return nil, 0, err
		}

		order.Customer = &customer

		items, err := r.getOrderItems(ctx, order.ID)
		if err != nil {
			return nil, 0, err
		}
		order.OrderItems = items

		orders = append(orders, order)
	}

	return orders, totalCount, nil
}

func (r *OrderRepo) GetOrdersByCustomer(ctx context.Context, userID string, customerID string, page, perPage uint) ([]models.Order, uint, error) {
	filter := models.OrderFilter{
		UserID:     userID,
		CustomerID: &customerID,
		Page:       page,
		PerPage:    perPage,
	}
	return r.GetOrders(ctx, filter)
}

func (r *OrderRepo) GetOrdersByCustomerOnly(ctx context.Context, filter models.OrderFilter) ([]models.Order, uint, error) {
	whereConditions := []string{}
	args := []interface{}{}
	argCount := 0

	if filter.CustomerID != nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("o.customer_id = $%d", argCount))
		args = append(args, *filter.CustomerID)
	} else {
		return nil, 0, fmt.Errorf("customer_id is required")
	}

	if filter.Status != nil {
		argCount++
		whereConditions = append(whereConditions, fmt.Sprintf("o.status = $%d", argCount))
		args = append(args, *filter.Status)
	}

	whereClause := strings.Join(whereConditions, " AND ")

	var totalCount uint
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM orders o WHERE %s`, whereClause)
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		log.Printf("[ERROR] Failed to get total count: %s", err.Error())
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PerPage
	argCount++
	limitArg := argCount
	argCount++
	offsetArg := argCount

	query := fmt.Sprintf(`
		SELECT 
			o.id, o.user_id, o.customer_id, o.total_price, o.status, o.order_date, o.created_at,
			c.id, c.name, c.address, c.phone, c.created_at
		FROM orders o
		LEFT JOIN customers c ON o.customer_id = c.id
		WHERE %s
		ORDER BY o.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, limitArg, offsetArg)

	args = append(args, filter.PerPage, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("[ERROR] Failed to get orders: %s", err.Error())
		return nil, 0, err
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var order models.Order
		var customer models.Customer
		err := rows.Scan(
			&order.ID, &order.UserID, &order.CustomerID, &order.TotalPrice,
			&order.Status, &order.OrderDate, &order.CreatedAt,
			&customer.ID, &customer.Name, &customer.Address, &customer.Phone, &customer.CreatedAt,
		)
		if err != nil {
			log.Printf("[ERROR] Failed to scan order: %s", err.Error())
			return nil, 0, err
		}

		order.Customer = &customer

		items, err := r.getOrderItems(ctx, order.ID)
		if err != nil {
			return nil, 0, err
		}
		order.OrderItems = items

		orders = append(orders, order)
	}

	return orders, totalCount, nil
}

func (r *OrderRepo) UpdateOrderStatus(ctx context.Context, orderID string, newStatus models.OrderStatus) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to begin transaction: %s", err.Error())
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("[ERROR] Failed to rollback transaction: %s", rbErr.Error())
			}
		}
	}()

	var currentStatus models.OrderStatus
	getOrderQuery := `
		SELECT status 
		FROM orders 
		WHERE id = $1 
		FOR UPDATE
	`
	err = tx.QueryRowContext(ctx, getOrderQuery, orderID).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("order not found")
		}
		log.Printf("[ERROR] Failed to get order: %s", err.Error())
		return err
	}

	if currentStatus == models.OrderStatusCancelled {
		return fmt.Errorf("cannot update status of cancelled order")
	}

	if currentStatus == models.OrderStatusConfirmed && newStatus == models.OrderStatusPending {
		return fmt.Errorf("cannot revert confirmed order to pending")
	}

	updateStatusQuery := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, updateStatusQuery, newStatus, orderID)
	if err != nil {
		log.Printf("[ERROR] Failed to update order status: %s", err.Error())
		return err
	}

	if newStatus == models.OrderStatusCancelled && currentStatus != models.OrderStatusCancelled {
		getItemsQuery := `
			SELECT product_id, quantity 
			FROM order_items 
			WHERE order_id = $1
		`
		rows, err := tx.QueryContext(ctx, getItemsQuery, orderID)
		if err != nil {
			log.Printf("[ERROR] Failed to get order items: %s", err.Error())
			return err
		}
		defer rows.Close()

		restockQty := make(map[string]int)

		for rows.Next() {
			var productID string
			var quantity int
			if err := rows.Scan(&productID, &quantity); err != nil {
				log.Printf("[ERROR] Failed to scan order item: %s", err.Error())
				return err
			}

			restockQty[productID] = quantity
		}

		for productID := range restockQty {
			restoreStockQuery := `
			UPDATE products 
			SET stock = stock + $1 
			WHERE id = $2
			`

			_, err = tx.ExecContext(ctx, restoreStockQuery, restockQty[productID], productID)
			if err != nil {
				log.Printf("[ERROR] Failed to restore stock: %s", err.Error())
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[ERROR] Failed to commit transaction: %s", err.Error())
		return err
	}

	return nil
}

func (r *OrderRepo) DeleteOrder(ctx context.Context, orderID string) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		log.Printf("[ERROR] Failed to begin transaction: %s", err.Error())
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("[ERROR] Failed to rollback transaction: %s", rbErr.Error())
			}
		}
	}()

	var status models.OrderStatus
	getOrderQuery := `
		SELECT status 
		FROM orders 
		WHERE id = $1 
		FOR UPDATE
	`
	err = tx.QueryRowContext(ctx, getOrderQuery, orderID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("order not found")
		}
		log.Printf("[ERROR] Failed to get order: %s", err.Error())
		return err
	}

	if status != models.OrderStatusPending && status != models.OrderStatusCancelled {
		return fmt.Errorf("can only delete pending or cancelled orders")
	}

	if status == models.OrderStatusPending {
		getItemsQuery := `
			SELECT product_id, quantity 
			FROM order_items 
			WHERE order_id = $1
		`
		rows, err := tx.QueryContext(ctx, getItemsQuery, orderID)
		if err != nil {
			log.Printf("[ERROR] Failed to get order items: %s", err.Error())
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var productID string
			var quantity int
			if err := rows.Scan(&productID, &quantity); err != nil {
				log.Printf("[ERROR] Failed to scan order item: %s", err.Error())
				return err
			}

			restoreStockQuery := `
				UPDATE products 
				SET stock = stock + $1 
				WHERE id = $2
			`
			_, err = tx.ExecContext(ctx, restoreStockQuery, quantity, productID)
			if err != nil {
				log.Printf("[ERROR] Failed to restore stock: %s", err.Error())
				return err
			}
		}
	}

	deleteItemsQuery := `DELETE FROM order_items WHERE order_id = $1`
	_, err = tx.ExecContext(ctx, deleteItemsQuery, orderID)
	if err != nil {
		log.Printf("[ERROR] Failed to delete order items: %s", err.Error())
		return err
	}

	deleteOrderQuery := `DELETE FROM orders WHERE id = $1`
	_, err = tx.ExecContext(ctx, deleteOrderQuery, orderID)
	if err != nil {
		log.Printf("[ERROR] Failed to delete order: %s", err.Error())
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Printf("[ERROR] Failed to commit transaction: %s", err.Error())
		return err
	}

	return nil
}

func (r *OrderRepo) GetCustomerByID(ctx context.Context, customerID int) (*models.Customer, error) {
	query := `
		SELECT id, name, address, phone, created_at
		FROM customers
		WHERE id = $1
	`

	var customer models.Customer
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(
		&customer.ID, &customer.Name, &customer.Address, &customer.Phone, &customer.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("customer not found")
		}
		log.Printf("[ERROR] Failed to get customer: %s", err.Error())
		return nil, err
	}

	return &customer, nil
}
