DROP TABLE IF EXISTS menu_items;
DROP TABLE IF EXISTS reviews;
DROP INDEX IF EXISTS idx_reviews_store_id;
DROP INDEX IF EXISTS idx_reviews_user_id;
DROP INDEX IF EXISTS idx_reviews_rating;
DROP INDEX IF EXISTS idx_reviews_visit_date;
DROP INDEX IF EXISTS idx_menu_items_review_id;
DROP INDEX IF EXISTS idx_reviews_unique_store_user;