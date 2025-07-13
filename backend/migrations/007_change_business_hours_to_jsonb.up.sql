-- Convert business_hours from TEXT to JSONB
-- First, create a temporary column
ALTER TABLE stores ADD COLUMN business_hours_jsonb JSONB DEFAULT '{
  "monday": {"is_closed": false, "time_slots": []},
  "tuesday": {"is_closed": false, "time_slots": []},
  "wednesday": {"is_closed": false, "time_slots": []},
  "thursday": {"is_closed": false, "time_slots": []},
  "friday": {"is_closed": false, "time_slots": []},
  "saturday": {"is_closed": false, "time_slots": []},
  "sunday": {"is_closed": false, "time_slots": []}
}';

-- Convert existing TEXT data to JSONB where possible
UPDATE stores 
SET business_hours_jsonb = CASE 
  WHEN business_hours IS NULL OR business_hours = '' THEN '{
    "monday": {"is_closed": false, "time_slots": []},
    "tuesday": {"is_closed": false, "time_slots": []},
    "wednesday": {"is_closed": false, "time_slots": []},
    "thursday": {"is_closed": false, "time_slots": []},
    "friday": {"is_closed": false, "time_slots": []},
    "saturday": {"is_closed": false, "time_slots": []},
    "sunday": {"is_closed": false, "time_slots": []}
  }'::jsonb
  WHEN business_hours::text ~ '^{.*}$' THEN business_hours::jsonb
  ELSE '{
    "monday": {"is_closed": false, "time_slots": []},
    "tuesday": {"is_closed": false, "time_slots": []},
    "wednesday": {"is_closed": false, "time_slots": []},
    "thursday": {"is_closed": false, "time_slots": []},
    "friday": {"is_closed": false, "time_slots": []},
    "saturday": {"is_closed": false, "time_slots": []},
    "sunday": {"is_closed": false, "time_slots": []}
  }'::jsonb
END;

-- Drop the old column
ALTER TABLE stores DROP COLUMN business_hours;

-- Rename the new column
ALTER TABLE stores RENAME COLUMN business_hours_jsonb TO business_hours;

-- Add index for JSONB queries
CREATE INDEX idx_stores_business_hours ON stores USING GIN (business_hours);