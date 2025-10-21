-- Create retrospective_groups table
CREATE TABLE retrospective_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    retrospective_id UUID NOT NULL REFERENCES retrospectives(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create retrospective_group_items table (junction table)
CREATE TABLE retrospective_group_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id UUID NOT NULL REFERENCES retrospective_groups(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES retrospective_items(id) ON DELETE CASCADE,
    UNIQUE(group_id, item_id)
);

-- Create retrospective_group_votes table
CREATE TABLE retrospective_group_votes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id UUID NOT NULL REFERENCES retrospective_groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(group_id, user_id)
);

-- Add vote count to retrospective_groups
ALTER TABLE retrospective_groups ADD COLUMN votes INTEGER NOT NULL DEFAULT 0;

-- Create indexes for better performance
CREATE INDEX idx_retrospective_groups_retrospective_id ON retrospective_groups(retrospective_id);
CREATE INDEX idx_retrospective_group_items_group_id ON retrospective_group_items(group_id);
CREATE INDEX idx_retrospective_group_items_item_id ON retrospective_group_items(item_id);
CREATE INDEX idx_retrospective_group_votes_group_id ON retrospective_group_votes(group_id);
CREATE INDEX idx_retrospective_group_votes_user_id ON retrospective_group_votes(user_id);
