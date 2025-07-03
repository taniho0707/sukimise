-- Database initialization script
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'editor',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Stores table
CREATE TABLE stores (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    categories JSONB DEFAULT '[]',
    business_hours TEXT,
    price_range VARCHAR(50),
    parking_info TEXT,
    website_url TEXT,
    google_map_url TEXT,
    sns_urls JSONB DEFAULT '[]',
    tags JSONB DEFAULT '[]',
    photos JSONB DEFAULT '[]',
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_stores_name ON stores(name);
CREATE INDEX idx_stores_location ON stores(latitude, longitude);
CREATE INDEX idx_stores_categories ON stores USING GIN (categories);
CREATE INDEX idx_stores_tags ON stores USING GIN (tags);
CREATE INDEX idx_stores_created_by ON stores(created_by);

-- Reviews table
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    photos JSONB DEFAULT '[]',
    visit_date TIMESTAMP WITH TIME ZONE,
    is_visited BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE menu_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    review_id UUID NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    comment TEXT
);

CREATE INDEX idx_reviews_store_id ON reviews(store_id);
CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_rating ON reviews(rating);
CREATE INDEX idx_reviews_visit_date ON reviews(visit_date);
CREATE INDEX idx_menu_items_review_id ON menu_items(review_id);

-- Removed unique constraint to allow multiple reviews per user per store
-- CREATE UNIQUE INDEX idx_reviews_unique_store_user ON reviews(store_id, user_id);

-- Insert a default admin user (password: admin123)
INSERT INTO users (username, email, password, role) VALUES (
    'admin',
    'admin@sukimise.com',
    '$2a$10$v2zOcygvW3kFIAWDVzsEeeQmTE0.dMWOtL7A1qr9eyRwTNMzWKdZG', -- bcrypt hash of 'admin123'
    'admin'
);

-- Insert a default editor user (password: editor123)
INSERT INTO users (username, email, password, role) VALUES (
    'editor',
    'editor@sukimise.com',
    '$2a$10$3dow5bs6VqqKAfYD2QwMieZYdLCime.DU5wTEccmtpTmopeo9upNC', -- bcrypt hash of 'editor123'
    'editor'
);