CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user', 
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE drivers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) UNIQUE,
    car_model VARCHAR(100),
    license_plate VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES users(id),
    driver_id INTEGER REFERENCES drivers(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('created', 'in_progress', 'done')),
    origin VARCHAR(255),
    destination VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE driver_locations (
    id SERIAL PRIMARY KEY,
    driver_id INTEGER REFERENCES drivers(id),
    order_id INTEGER REFERENCES orders(id),
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    generated_time TIMESTAMP NOT NULL,     
    recorded_time TIMESTAMP DEFAULT NOW()  
);