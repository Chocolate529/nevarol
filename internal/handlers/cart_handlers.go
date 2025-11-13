package handlers

import (
	"net/http"
	"strconv"

	"github.com/Chocolate529/nevarol/internal/email"
	"github.com/Chocolate529/nevarol/internal/models"
	"github.com/go-chi/chi/v5"
)

// GetCart returns the user's cart items
func (m *Repository) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := m.App.Session.GetInt(r.Context(), "user_id")
	if userID == 0 {
		writeJSON(w, http.StatusUnauthorized, JSONResponse{
			OK:      false,
			Message: "Not authenticated",
		})
		return
	}

	items, err := m.App.DB.GetCartItems(userID)
	if err != nil {
		m.App.ErrorLog.Println("Error getting cart items:", err)
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			OK:      false,
			Message: "Failed to get cart items",
		})
		return
	}

	// Return empty array instead of null
	if items == nil {
		items = []models.CartItem{}
	}

	writeJSON(w, http.StatusOK, JSONResponse{
		OK:   true,
		Data: items,
	})
}

// AddToCart adds an item to the cart
func (m *Repository) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := m.App.Session.GetInt(r.Context(), "user_id")
	if userID == 0 {
		writeJSON(w, http.StatusUnauthorized, JSONResponse{
			OK:      false,
			Message: "Not authenticated",
		})
		return
	}

	var payload struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	err := readJSON(w, r, &payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Invalid request format",
		})
		return
	}

	if payload.ProductID <= 0 || payload.Quantity <= 0 {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Invalid product ID or quantity",
		})
		return
	}

	err = m.App.DB.AddToCart(userID, payload.ProductID, payload.Quantity)
	if err != nil {
		m.App.ErrorLog.Println("Error adding to cart:", err)
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			OK:      false,
			Message: "Failed to add to cart",
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{
		OK:      true,
		Message: "Item added to cart",
	})
}

// UpdateCartItem updates the quantity of a cart item
func (m *Repository) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userID := m.App.Session.GetInt(r.Context(), "user_id")
	if userID == 0 {
		writeJSON(w, http.StatusUnauthorized, JSONResponse{
			OK:      false,
			Message: "Not authenticated",
		})
		return
	}

	itemIDStr := chi.URLParam(r, "id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Invalid item ID",
		})
		return
	}

	var payload struct {
		Quantity int `json:"quantity"`
	}

	err = readJSON(w, r, &payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Invalid request format",
		})
		return
	}

	if payload.Quantity <= 0 {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Quantity must be greater than 0",
		})
		return
	}

	err = m.App.DB.UpdateCartItem(itemID, payload.Quantity)
	if err != nil {
		m.App.ErrorLog.Println("Error updating cart item:", err)
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			OK:      false,
			Message: "Failed to update cart item",
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{
		OK:      true,
		Message: "Cart item updated",
	})
}

// RemoveFromCart removes an item from the cart
func (m *Repository) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := m.App.Session.GetInt(r.Context(), "user_id")
	if userID == 0 {
		writeJSON(w, http.StatusUnauthorized, JSONResponse{
			OK:      false,
			Message: "Not authenticated",
		})
		return
	}

	itemIDStr := chi.URLParam(r, "id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Invalid item ID",
		})
		return
	}

	err = m.App.DB.RemoveFromCart(itemID)
	if err != nil {
		m.App.ErrorLog.Println("Error removing from cart:", err)
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			OK:      false,
			Message: "Failed to remove from cart",
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{
		OK:      true,
		Message: "Item removed from cart",
	})
}

// ClearCart removes all items from the user's cart
func (m *Repository) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID := m.App.Session.GetInt(r.Context(), "user_id")
	if userID == 0 {
		writeJSON(w, http.StatusUnauthorized, JSONResponse{
			OK:      false,
			Message: "Not authenticated",
		})
		return
	}

	err := m.App.DB.ClearCart(userID)
	if err != nil {
		m.App.ErrorLog.Println("Error clearing cart:", err)
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			OK:      false,
			Message: "Failed to clear cart",
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{
		OK:      true,
		Message: "Cart cleared",
	})
}

// CreateOrder creates an order from the user's cart
func (m *Repository) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := m.App.Session.GetInt(r.Context(), "user_id")
	if userID == 0 {
		writeJSON(w, http.StatusUnauthorized, JSONResponse{
			OK:      false,
			Message: "Not authenticated",
		})
		return
	}

	// Parse contact information from request
	var payload struct {
		CustomerName  string `json:"customer_name"`
		CustomerEmail string `json:"customer_email"`
		Phone         string `json:"phone"`
		Address       string `json:"address"`
	}

	err := readJSON(w, r, &payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate contact information
	if payload.CustomerName == "" || payload.CustomerEmail == "" || 
	   payload.Phone == "" || payload.Address == "" {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "All contact fields are required (name, email, phone, address)",
		})
		return
	}

	order, err := m.App.DB.CreateOrder(userID, payload.CustomerName, payload.CustomerEmail, payload.Phone, payload.Address)
	if err != nil {
		m.App.ErrorLog.Println("Error creating order:", err)
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			OK:      false,
			Message: "Failed to create order",
		})
		return
	}

	// Send email notifications
	if m.App.EmailConfig.IsConfigured() {
		// Prepare order details for email
		var emailItems []email.OrderItemDetail
		for _, item := range order.Items {
			emailItems = append(emailItems, email.OrderItemDetail{
				ProductName: item.Product.Name,
				Quantity:    item.Quantity,
				Price:       item.Price,
			})
		}

		orderDetails := email.OrderDetails{
			OrderID:       order.ID,
			CustomerEmail: order.CustomerEmail,
			CustomerName:  order.CustomerName,
			Phone:         order.Phone,
			Address:       order.Address,
			TotalPrice:    order.TotalPrice,
			Items:         emailItems,
		}

		// Send notification to admin
		err = m.App.EmailConfig.SendOrderNotification(orderDetails)
		if err != nil {
			m.App.ErrorLog.Println("Failed to send admin notification:", err)
			// Don't fail the order if email fails
		}

		// Send confirmation to customer
		err = m.App.EmailConfig.SendOrderConfirmation(orderDetails)
		if err != nil {
			m.App.ErrorLog.Println("Failed to send customer confirmation:", err)
			// Don't fail the order if email fails
		}
	}

	writeJSON(w, http.StatusCreated, JSONResponse{
		OK:      true,
		Message: "Order created successfully",
		Data:    order,
	})
}

// GetOrders returns the user's orders
func (m *Repository) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := m.App.Session.GetInt(r.Context(), "user_id")
	if userID == 0 {
		writeJSON(w, http.StatusUnauthorized, JSONResponse{
			OK:      false,
			Message: "Not authenticated",
		})
		return
	}

	orders, err := m.App.DB.GetUserOrders(userID)
	if err != nil {
		m.App.ErrorLog.Println("Error getting orders:", err)
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			OK:      false,
			Message: "Failed to get orders",
		})
		return
	}

	// Return empty array instead of null
	if orders == nil {
		orders = []models.Order{}
	}

	writeJSON(w, http.StatusOK, JSONResponse{
		OK:   true,
		Data: orders,
	})
}
