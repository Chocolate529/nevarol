package repository

import (
	"context"
	"time"

	"github.com/Chocolate529/nevarol/internal/models"
)

// GetAllProducts retrieves all products
func (m *DatabaseRepo) GetAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, name, price, type, image, description FROM products ORDER BY id`

	rows, err := m.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Type, &p.Image, &p.Description)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// GetProductByID retrieves a product by ID
func (m *DatabaseRepo) GetProductByID(id int) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var product models.Product
	query := `SELECT id, name, price, type, image, description FROM products WHERE id = $1`

	err := m.DB.QueryRow(ctx, query, id).Scan(
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

	return &product, nil
}
