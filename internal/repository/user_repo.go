package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Chocolate529/nevarol/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type DatabaseRepo struct {
	DB *pgxpool.Pool
}

// NewDatabaseRepo creates a new database repository
func NewDatabaseRepo(db *pgxpool.Pool) *DatabaseRepo {
	return &DatabaseRepo{
		DB: db,
	}
}

// CreateUser creates a new user with hashed password
func (m *DatabaseRepo) CreateUser(email, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	var user models.User
	query := `
		INSERT INTO users (email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, created_at, updated_at
	`

	now := time.Now()
	err = m.DB.QueryRow(ctx, query, email, string(hashedPassword), now, now).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (m *DatabaseRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`

	err := m.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// AuthenticateUser validates user credentials
func (m *DatabaseRepo) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := m.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Clear password before returning
	user.Password = ""
	return user, nil
}

// GetUserByID retrieves a user by ID
func (m *DatabaseRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	query := `SELECT id, email, created_at, updated_at FROM users WHERE id = $1`

	err := m.DB.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
