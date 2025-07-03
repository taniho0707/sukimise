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