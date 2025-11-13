package repository

import (
	"context"
	"time"

	"github.com/Chocolate529/nevarol/internal/models"
)

// GetCartItems retrieves all cart items for a user
func (m *DatabaseRepo) GetCartItems(userID int) ([]models.CartItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT c.id, c.user_id, c.product_id, c.quantity,
		       p.id, p.name, p.price, p.type, p.image, p.description
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = $1
	`

	rows, err := m.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		var product models.Product
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ProductID,
			&item.Quantity,
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Type,
			&product.Image,
			&product.Description,
		)
		if err != nil {
			return nil, err
		}
		item.Product = product
		items = append(items, item)
	}

	return items, nil
}

// AddToCart adds an item to the cart or updates quantity if it exists
func (m *DatabaseRepo) AddToCart(userID, productID, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if item already exists in cart
	var existingID int
	var existingQty int
	checkQuery := `SELECT id, quantity FROM cart_items WHERE user_id = $1 AND product_id = $2`
	err := m.DB.QueryRow(ctx, checkQuery, userID, productID).Scan(&existingID, &existingQty)

	if err == nil {
		// Item exists, update quantity
		updateQuery := `UPDATE cart_items SET quantity = $1 WHERE id = $2`
		_, err = m.DB.Exec(ctx, updateQuery, existingQty+quantity, existingID)
		return err
	}

	// Item doesn't exist, insert new
	insertQuery := `INSERT INTO cart_items (user_id, product_id, quantity) VALUES ($1, $2, $3)`
	_, err = m.DB.Exec(ctx, insertQuery, userID, productID, quantity)
	return err
}

// UpdateCartItem updates the quantity of a cart item
func (m *DatabaseRepo) UpdateCartItem(itemID, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE cart_items SET quantity = $1 WHERE id = $2`
	_, err := m.DB.Exec(ctx, query, quantity, itemID)
	return err
}

// RemoveFromCart removes an item from the cart
func (m *DatabaseRepo) RemoveFromCart(itemID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM cart_items WHERE id = $1`
	_, err := m.DB.Exec(ctx, query, itemID)
	return err
}

// ClearCart removes all items from a user's cart
func (m *DatabaseRepo) ClearCart(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM cart_items WHERE user_id = $1`
	_, err := m.DB.Exec(ctx, query, userID)
	return err
}

// CreateOrder creates a new order from cart items
func (m *DatabaseRepo) CreateOrder(userID int) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start transaction
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Get cart items
	cartQuery := `
		SELECT c.product_id, c.quantity, p.price
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = $1
	`
	rows, err := tx.Query(ctx, cartQuery, userID)
	if err != nil {
		return nil, err
	}

	var totalPrice float64
	var orderItems []struct {
		ProductID int
		Quantity  int
		Price     float64
	}

	for rows.Next() {
		var item struct {
			ProductID int
			Quantity  int
			Price     float64
		}
		err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price)
		if err != nil {
			rows.Close()
			return nil, err
		}
		totalPrice += item.Price * float64(item.Quantity)
		orderItems = append(orderItems, item)
	}
	rows.Close()

	// Create order
	var orderID int
	orderQuery := `
		INSERT INTO orders (user_id, total_price, status, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err = tx.QueryRow(ctx, orderQuery, userID, totalPrice, "pending", time.Now()).Scan(&orderID)
	if err != nil {
		return nil, err
	}

	// Create order items
	for _, item := range orderItems {
		itemQuery := `
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`
		_, err = tx.Exec(ctx, itemQuery, orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return nil, err
		}
	}

	// Clear cart
	clearQuery := `DELETE FROM cart_items WHERE user_id = $1`
	_, err = tx.Exec(ctx, clearQuery, userID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &models.Order{
		ID:         orderID,
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "pending",
	}, nil
}

// GetUserOrders retrieves all orders for a user
func (m *DatabaseRepo) GetUserOrders(userID int) ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, total_price, status, created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := m.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
