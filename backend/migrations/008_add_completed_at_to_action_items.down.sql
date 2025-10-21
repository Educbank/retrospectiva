-- Remove completed_at field from action_items table
DROP INDEX IF EXISTS idx_action_items_completed_at;
ALTER TABLE action_items DROP COLUMN IF EXISTS completed_at;
