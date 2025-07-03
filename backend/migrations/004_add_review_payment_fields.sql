-- Add payment amount and food notes to reviews table
ALTER TABLE reviews ADD COLUMN payment_amount INTEGER; -- 支払金額（円）
ALTER TABLE reviews ADD COLUMN food_notes TEXT; -- 料理についてのメモ

-- Add index for payment amount for efficient average calculation
CREATE INDEX idx_reviews_payment_amount ON reviews(store_id, payment_amount) WHERE payment_amount IS NOT NULL;