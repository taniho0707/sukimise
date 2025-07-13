-- Rollback: Convert business_hours from JSONB back to TEXT
-- Drop the GIN index first
DROP INDEX IF EXISTS idx_stores_business_hours;

-- Add a temporary TEXT column
ALTER TABLE stores ADD COLUMN business_hours_text TEXT;

-- Convert JSONB back to TEXT (simplified format)
UPDATE stores 
SET business_hours_text = COALESCE(business_hours::text, '');

-- Drop the JSONB column
ALTER TABLE stores DROP COLUMN business_hours;

-- Rename the TEXT column back
ALTER TABLE stores RENAME COLUMN business_hours_text TO business_hours;