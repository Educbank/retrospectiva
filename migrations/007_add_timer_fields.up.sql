-- Add timer fields to retrospectives table
ALTER TABLE retrospectives 
ADD COLUMN timer_duration INTEGER DEFAULT 0, -- Timer duration in seconds
ADD COLUMN timer_started_at TIMESTAMP NULL, -- When timer was started
ADD COLUMN timer_paused_at TIMESTAMP NULL, -- When timer was paused
ADD COLUMN timer_elapsed_time INTEGER DEFAULT 0; -- Total elapsed time in seconds
