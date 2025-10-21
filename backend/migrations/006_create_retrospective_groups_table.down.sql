-- Drop indexes
DROP INDEX IF EXISTS idx_retrospective_group_votes_user_id;
DROP INDEX IF EXISTS idx_retrospective_group_votes_group_id;
DROP INDEX IF EXISTS idx_retrospective_group_items_item_id;
DROP INDEX IF EXISTS idx_retrospective_group_items_group_id;
DROP INDEX IF EXISTS idx_retrospective_groups_retrospective_id;

-- Drop tables
DROP TABLE IF EXISTS retrospective_group_votes;
DROP TABLE IF EXISTS retrospective_group_items;
DROP TABLE IF EXISTS retrospective_groups;
