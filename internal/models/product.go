package models

// Product represents a product in the store
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Type        string  `json:"type"`
	Image       string  `json:"image"`
	Description string  `json:"description"`
}
