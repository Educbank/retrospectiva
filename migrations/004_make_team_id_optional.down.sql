-- Revert team_id to required in retrospectives table
ALTER TABLE retrospectives ALTER COLUMN team_id SET NOT NULL;
