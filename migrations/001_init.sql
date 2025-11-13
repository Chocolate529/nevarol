-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    type VARCHAR(50) NOT NULL,
    image VARCHAR(255),
    description TEXT
);

-- Cart items table
CREATE TABLE IF NOT EXISTS cart_items (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    UNIQUE(user_id, product_id)
);

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    address TEXT NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Order items table
CREATE TABLE IF NOT EXISTS order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL
);

-- Insert initial products
INSERT INTO products (name, price, type, image, description) VALUES
    ('Polyurethane Wheel Ø80mm', 19.99, 'polyurethane', 'images/wheel1.jpg', 'High-quality polyurethane wheel, 80mm diameter'),
    ('Nylon Wheel Ø70mm', 14.50, 'nylon', 'images/wheel2.jpg', 'Durable nylon wheel, 70mm diameter'),
    ('Rubber Coated Wheel Ø90mm', 22.00, 'rubber', 'images/wheel3.jpg', 'Rubber coated wheel for smooth operation, 90mm'),
    ('Polyurethane Wheel Ø100mm', 25.00, 'polyurethane', 'images/wheel4.jpg', 'Premium polyurethane wheel, 100mm diameter'),
    ('Nylon Wheel Ø85mm', 17.80, 'nylon', 'images/wheel5.jpg', 'Heavy-duty nylon wheel, 85mm'),
    ('Rubber Wheel Ø75mm', 16.20, 'rubber', 'images/wheel6.jpg', 'Shock-absorbing rubber wheel, 75mm'),
    ('Polyurethane Wheel Ø110mm', 28.40, 'polyurethane', 'images/wheel7.jpg', 'Extra large polyurethane wheel, 110mm'),
    ('Nylon Heavy Duty Ø95mm', 20.00, 'nylon', 'images/wheel8.jpg', 'Industrial-grade nylon wheel, 95mm'),
    ('Rubber Shock-Absorb Ø100mm', 27.50, 'rubber', 'images/wheel9.jpg', 'Premium shock-absorbing rubber wheel, 100mm'),
    ('Polyurethane Silent Ø90mm', 23.90, 'polyurethane', 'images/wheel10.jpg', 'Silent operation polyurethane wheel, 90mm')
ON CONFLICT DO NOTHING;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
