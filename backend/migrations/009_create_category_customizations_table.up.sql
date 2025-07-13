-- Drop old tag customizations table
DROP TABLE IF EXISTS tag_customizations;

-- Add category customizations table for custom icons and colors
CREATE TABLE category_customizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_name VARCHAR(255) UNIQUE NOT NULL,
    icon VARCHAR(10), -- Single emoji or character
    color VARCHAR(7), -- Hex color code like #FF5733
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create index for faster category lookups
CREATE INDEX idx_category_customizations_category_name ON category_customizations(category_name);

-- Insert default customizations for common categories
INSERT INTO category_customizations (category_name, icon, color) VALUES
('レストラン', '🍽️', '#FF5733'),
('カフェ', '☕', '#8B4513'),
('ラーメン', '🍜', '#FF6B35'),
('寿司', '🍣', '#FF1744'),
('居酒屋', '🍻', '#FFA726'),
('ファストフード', '🍔', '#FF9800'),
('イタリアン', '🍝', '#E91E63'),
('中華', '🥟', '#F44336'),
('焼肉', '🥩', '#795548'),
('和食', '🍱', '#4CAF50'),
('パン屋', '🍞', '#FF7043'),
('スイーツ', '🧁', '#E1BEE7'),
('コンビニ', '🏪', '#2196F3'),
('スーパー', '🛒', '#009688'),
('病院', '🏥', '#F44336'),
('薬局', '💊', '#9C27B0'),
('銀行', '🏦', '#607D8B'),
('郵便局', '📮', '#FF5722'),
('駐車場', '🅿️', '#795548'),
('ガソリンスタンド', '⛽', '#FF9800'),
('美容院', '💇', '#E91E63'),
('書店', '📚', '#3F51B5');