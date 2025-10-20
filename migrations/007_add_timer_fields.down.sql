-- Remove timer fields from retrospectives table
ALTER TABLE retrospectives 
DROP COLUMN timer_duration,
DROP COLUMN timer_started_at,
DROP COLUMN timer_paused_at,
DROP COLUMN timer_elapsed_time;
