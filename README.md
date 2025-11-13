# Nevarol - E-commerce Web Application

A modern e-commerce web application for selling pallet truck wheels, built with Go and PostgreSQL.

## Features

- **Secure Authentication**: User registration and login with bcrypt password hashing
- **Database Persistence**: All data stored in PostgreSQL database
- **Product Management**: Dynamic product catalog with filtering and search
- **Shopping Cart**: Persistent cart with real-time updates
- **Order Management**: Complete order processing and history
- **Session Management**: Secure session handling with CSRF protection
- **Modern UI**: Responsive Bootstrap 5 interface

## Tech Stack

- **Backend**: Go 1.24.5
- **Database**: PostgreSQL
- **Router**: Chi v5
- **Session**: SCS v2
- **Frontend**: Bootstrap 5, Vanilla JavaScript
- **Security**: bcrypt, CSRF protection (nosurf)

## Prerequisites

- Go 1.24.5 or later
- PostgreSQL 12 or later

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Chocolate529/nevarol.git
cd nevarol
```

2. Install Go dependencies:
```bash
go mod download
```

3. Create a PostgreSQL database:
```bash
createdb nevarol
```

4. Set up environment variables:
```bash
cp .env.example .env
```

Edit `.env` with your database credentials:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=nevarol
IN_PRODUCTION=false
```

5. Build and run the application:
```bash
go build -o app ./cmd/web/
./app
```

The application will automatically run database migrations on startup.

## Usage

1. Open your browser and navigate to `http://localhost:8080`
2. Register a new account or login
3. Browse products in the store
4. Add items to your cart
5. Proceed to checkout

## Database Schema

The application uses the following tables:
- `users`: User accounts with hashed passwords
- `products`: Product catalog
- `cart_items`: Shopping cart items
- `orders`: Completed orders
- `order_items`: Order line items

## Security Features

- Password hashing with bcrypt (cost factor 12)
- CSRF protection on all state-changing requests
- Secure session cookies
- Input validation and sanitization
- SQL injection protection via parameterized queries
- XSS protection

## Development

To run in development mode:
```bash
go run ./cmd/web/
```

## License

This project is licensed under the MIT License.

