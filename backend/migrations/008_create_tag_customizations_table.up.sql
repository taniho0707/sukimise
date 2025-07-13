-- Add tag customizations table for custom icons and colors
CREATE TABLE tag_customizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tag_name VARCHAR(255) UNIQUE NOT NULL,
    icon VARCHAR(10), -- Single emoji or character
    color VARCHAR(7), -- Hex color code like #FF5733
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create index for faster tag lookups
CREATE INDEX idx_tag_customizations_tag_name ON tag_customizations(tag_name);

-- Insert default customizations for common tags
INSERT INTO tag_customizations (tag_name, icon, color) VALUES
('レストラン', '🍽️', '#FF5733'),
('カフェ', '☕', '#8B4513'),
('ラーメン', '🍜', '#FF6B35'),
('寿司', '🍣', '#FF1744'),
('居酒屋', '🍻', '#FFA726'),
('コンビニ', '🏪', '#4CAF50'),
('スーパー', '🛒', '#2196F3'),
('病院', '🏥', '#F44336'),
('薬局', '💊', '#9C27B0'),
('銀行', '🏦', '#607D8B'),
('郵便局', '📮', '#FF5722'),
('駐車場', '🅿️', '#795548');