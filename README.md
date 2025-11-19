# Nevarol - E-commerce Web Application

A modern e-commerce web application for selling pallet truck wheels, built with Go and PostgreSQL.

## 🎯 New to This Project?

**Wondering about architecture and hosting?** See our comprehensive guides:

- 📖 **[ARCHITECTURE_RECOMMENDATIONS.md](ARCHITECTURE_RECOMMENDATIONS.md)** - Should you use single-page or multi-page? What hosting to choose?
- 🚀 **[HOSTING_GUIDE.md](HOSTING_GUIDE.md)** - Step-by-step deployment to Fly.io, Railway, or DigitalOcean
- 📧 **[EMAIL_SETUP.md](EMAIL_SETUP.md)** - Configure order notifications
- ✅ **[PRODUCTION_READINESS.md](PRODUCTION_READINESS.md)** - Pre-launch checklist

## Features

- **Secure Authentication**: User registration and login with bcrypt password hashing
- **Database Persistence**: All data stored in PostgreSQL database
- **Product Management**: Dynamic product catalog with filtering and search
- **Shopping Cart**: Persistent cart with real-time updates
- **Order Management**: Complete order processing and history with contact information
- **Email Notifications**: Automatic order notifications to admin and customers (optional)
- **Session Management**: Secure session handling with CSRF protection
- **Rate Limiting**: Built-in rate limiting to prevent abuse
- **Security Headers**: Comprehensive security headers (CSP, X-Frame-Options, etc.)
- **Modern UI**: Responsive Bootstrap 5 interface

## Tech Stack

- **Backend**: Go 1.24.5
- **Database**: PostgreSQL 15
- **Router**: Chi v5
- **Session**: SCS v2
- **Email**: Go's built-in SMTP (supports Gmail, SendGrid, Mailgun, AWS SES)
- **Frontend**: Bootstrap 5, Vanilla JavaScript
- **Security**: bcrypt, CSRF protection (nosurf), rate limiting

## Prerequisites

- Go 1.24.5 or later
- PostgreSQL 12 or later
- OR Docker and Docker Compose
- (Optional) Email account for order notifications

## Installation

### Option 1: Using Docker (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/Chocolate529/nevarol.git
cd nevarol
```

2. Start the application using Docker Compose:
```bash
docker-compose up -d
```

The application will be available at `http://localhost:8080`

To stop the application:
```bash
docker-compose down
```

### Option 2: Manual Installation

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

# Optional: Email notifications
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
ADMIN_EMAIL=your-admin-email@gmail.com
```

**Note**: Email configuration is optional. See `EMAIL_SETUP.md` for detailed email setup instructions.

5. Build and run the application:
```bash
go build -o app ./cmd/web/
./app
```

The application will automatically run database migrations on startup.

## Email Notifications (Optional)

The application can send email notifications when orders are placed:
- **Admin receives**: Order details with customer contact information
- **Customer receives**: Order confirmation with contact details

**Email is completely optional** - orders work without email configuration.

For setup instructions, see [EMAIL_SETUP.md](EMAIL_SETUP.md).

Quick Gmail setup:
1. Generate App Password at https://myaccount.google.com/apppasswords
2. Add to `.env`:
   ```env
   SMTP_USER=your-email@gmail.com
   SMTP_PASSWORD=your-16-char-app-password
   ADMIN_EMAIL=your-email@gmail.com
   ```

## Usage

1. Open your browser and navigate to `http://localhost:8080`
2. Register a new account or login
3. Browse products in the store
4. Add items to your cart
5. Click checkout and enter your contact information
6. Receive order confirmation (via email if configured)

## Database Schema

The application uses the following tables:
- `users`: User accounts with hashed passwords
- `products`: Product catalog (pre-populated with 10 wheel products)
- `cart_items`: Shopping cart items
- `orders`: Completed orders with customer contact information
- `order_items`: Order line items

## Security Features

- **Password Security**: bcrypt hashing with cost factor 12
- **CSRF Protection**: Enabled on all state-changing requests
- **Secure Cookies**: HttpOnly, SameSite, and Secure flags
- **Rate Limiting**: 100 requests per minute per IP
- **Security Headers**:
  - X-Frame-Options: DENY (prevent clickjacking)
  - X-Content-Type-Options: nosniff (prevent MIME sniffing)
  - X-XSS-Protection: 1; mode=block
  - Content-Security-Policy: Restrictive CSP
  - Referrer-Policy: strict-origin-when-cross-origin
- **Input Validation**: All user inputs validated and sanitized
- **SQL Injection Protection**: Parameterized queries via pgx
- **Session Security**: Secure session management with automatic cleanup

## API Endpoints

### Authentication
- `POST /api/register` - Register a new user
- `POST /api/login` - Login
- `POST /api/logout` - Logout
- `GET /api/user` - Get current user

### Products
- `GET /api/products` - Get all products

### Cart
- `GET /api/cart` - Get cart items
- `POST /api/cart` - Add item to cart
- `PUT /api/cart/{id}` - Update cart item quantity
- `DELETE /api/cart/{id}` - Remove item from cart
- `DELETE /api/cart` - Clear cart

### Orders
- `POST /api/orders` - Create order from cart
- `GET /api/orders` - Get user's orders

## Development

To run in development mode:
```bash
go run ./cmd/web/
```

To run tests:
```bash
go test ./...
```

### GoLand IDE with Docker

If you're using GoLand IDE, see [GOLAND_DOCKER_SETUP.md](GOLAND_DOCKER_SETUP.md) for a comprehensive guide on:
- Setting up GoLand with Docker
- Running and debugging the application
- Using database tools
- Common workflows and troubleshooting

## Production Deployment

For production deployment:

1. Set `IN_PRODUCTION=true` in your environment
2. Use HTTPS/TLS (required for secure cookies)
3. Use strong database credentials
4. Configure proper backup strategy for PostgreSQL
5. Consider using a reverse proxy (nginx/Caddy) for additional security
6. Enable database connection pooling settings as needed
7. Monitor rate limits and adjust as necessary

## License

This project is licensed under the MIT License.

