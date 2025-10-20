-- Make team_id optional in retrospectives table
ALTER TABLE retrospectives ALTER COLUMN team_id DROP NOT NULL;
