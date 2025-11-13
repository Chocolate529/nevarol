package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Chocolate529/nevarol/internal/models"
)

// JSONResponse is a standard JSON response structure
type JSONResponse struct {
	OK      bool        `json:"ok"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// readJSON reads JSON from request body
func readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	
	return nil
}

// Register handles user registration
func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := readJSON(w, r, &payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate input
	if payload.Email == "" || payload.Password == "" {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Email and password are required",
		})
		return
	}

	// Basic email validation
	if !strings.Contains(payload.Email, "@") {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Invalid email format",
		})
		return
	}

	// Password validation
	if len(payload.Password) < 6 {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Password must be at least 6 characters",
		})
		return
	}

	// Create user
	user, err := m.App.DB.CreateUser(payload.Email, payload.Password)
	if err != nil {
		// Check if user already exists (duplicate email)
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			writeJSON(w, http.StatusConflict, JSONResponse{
				OK:      false,
				Message: "Email already registered",
			})
			return
		}

		m.App.ErrorLog.Println("Error creating user:", err)
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			OK:      false,
			Message: "Failed to create user",
		})
		return
	}

	writeJSON(w, http.StatusCreated, JSONResponse{
		OK:      true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// LoginAPI handles user login
func (m *Repository) LoginAPI(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := readJSON(w, r, &payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate input
	if payload.Email == "" || payload.Password == "" {
		writeJSON(w, http.StatusBadRequest, JSONResponse{
			OK:      false,
			Message: "Email and password are required",
		})
		return
	}

	// Authenticate user
	user, err := m.App.DB.AuthenticateUser(payload.Email, payload.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, JSONResponse{
			OK:      false,
			Message: "Invalid credentials",
		})
		return
	}

	// Store user in session
	m.App.Session.Put(r.Context(), "user_id", user.ID)
	m.App.Session.Put(r.Context(), "user_email", user.Email)

	writeJSON(w, http.StatusOK, JSONResponse{
		OK:      true,
		Message: "Login successful",
		Data:    user,
	})
}

// LogoutAPI handles user logout
func (m *Repository) LogoutAPI(w http.ResponseWriter, r *http.Request) {
	// Clear session
	err := m.App.Session.Destroy(r.Context())
	if err != nil {
		m.App.ErrorLog.Println("Error destroying session:", err)
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			OK:      false,
			Message: "Failed to logout",
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{
		OK:      true,
		Message: "Logout successful",
	})
}

// GetCurrentUser returns the currently logged in user
func (m *Repository) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := m.App.Session.GetInt(r.Context(), "user_id")
	if userID == 0 {
		writeJSON(w, http.StatusUnauthorized, JSONResponse{
			OK:      false,
			Message: "Not authenticated",
		})
		return
	}

	user, err := m.App.DB.GetUserByID(userID)
	if err != nil {
		m.App.ErrorLog.Println("Error getting user:", err)
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			OK:      false,
			Message: "Failed to get user",
		})
		return
	}

	writeJSON(w, http.StatusOK, JSONResponse{
		OK:   true,
		Data: user,
	})
}

// GetProducts returns all products
func (m *Repository) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := m.App.DB.GetAllProducts()
	if err != nil {
		m.App.ErrorLog.Println("Error getting products:", err)
		writeJSON(w, http.StatusInternalServerError, JSONResponse{
			OK:      false,
			Message: "Failed to get products",
		})
		return
	}

	// If no products, return empty array instead of null
	if products == nil {
		products = []models.Product{}
	}

	writeJSON(w, http.StatusOK, JSONResponse{
		OK:   true,
		Data: products,
	})
}
