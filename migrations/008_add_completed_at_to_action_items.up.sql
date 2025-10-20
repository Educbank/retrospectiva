-- Add completed_at field to action_items table
ALTER TABLE action_items ADD COLUMN completed_at TIMESTAMP WITH TIME ZONE;

-- Create index for completed_at field
CREATE INDEX idx_action_items_completed_at ON action_items(completed_at);
