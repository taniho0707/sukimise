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

-- Default users are now created via environment variables and admin commands
-- See README.md for instructions on creating users

-- Create viewer password settings table
CREATE TABLE viewer_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    password_hash VARCHAR(255) NOT NULL,
    session_duration_days INTEGER NOT NULL DEFAULT 7,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create viewer login history table
CREATE TABLE viewer_login_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip_address INET,
    user_agent TEXT,
    login_time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    session_token VARCHAR(255),
    expires_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better performance
CREATE INDEX idx_viewer_login_history_login_time ON viewer_login_history(login_time);
CREATE INDEX idx_viewer_login_history_session_token ON viewer_login_history(session_token);
CREATE INDEX idx_viewer_login_history_expires_at ON viewer_login_history(expires_at);

-- Insert default viewer settings (password: viewer123)
INSERT INTO viewer_settings (password_hash, session_duration_days) VALUES (
    '$2a$10$vPZxOoHW8tRYvBhDHN4yBOmJQfgVzv7rVHvLFxEGIGsNTVcBjJqhS', -- bcrypt hash of 'viewer123'
    7
);