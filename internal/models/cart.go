package models

// CartItem represents an item in a user's shopping cart
type CartItem struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Product   Product `json:"product,omitempty"`
}

// Order represents a completed order
type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  string    `json:"created_at"`
	Items      []OrderItem `json:"items,omitempty"`
}

// OrderItem represents a single item in an order
type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Product   Product `json:"product,omitempty"`
}
