-- Revert template constraint to original values
ALTER TABLE retrospectives DROP CONSTRAINT IF EXISTS retrospectives_template_check;
ALTER TABLE retrospectives ADD CONSTRAINT retrospectives_template_check 
CHECK (template IN ('start_stop_continue', '4ls', 'mad_sad_glad', 'sailboat', 'custom'));
