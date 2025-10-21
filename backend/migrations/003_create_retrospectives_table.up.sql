CREATE TABLE retrospectives (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    template VARCHAR(50) NOT NULL CHECK (template IN ('start_stop_continue', '4ls', 'mad_sad_glad', 'sailboat', 'went_well_to_improve', 'custom')),
    status VARCHAR(50) NOT NULL DEFAULT 'planned' CHECK (status IN ('planned', 'active', 'collecting', 'voting', 'discussing', 'closed')),
    scheduled_at TIMESTAMP WITH TIME ZONE,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE retrospective_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    retrospective_id UUID NOT NULL REFERENCES retrospectives(id) ON DELETE CASCADE,
    category VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    author_id UUID REFERENCES users(id) ON DELETE SET NULL,
    is_anonymous BOOLEAN NOT NULL DEFAULT false,
    votes INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE retrospective_votes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id UUID NOT NULL REFERENCES retrospective_items(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(item_id, user_id)
);

CREATE TABLE action_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    retrospective_id UUID NOT NULL REFERENCES retrospectives(id) ON DELETE CASCADE,
    item_id UUID REFERENCES retrospective_items(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'todo' CHECK (status IN ('todo', 'in_progress', 'done')),
    due_date TIMESTAMP WITH TIME ZONE,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_retrospectives_team_id ON retrospectives(team_id);
CREATE INDEX idx_retrospectives_status ON retrospectives(status);
CREATE INDEX idx_retrospectives_created_by ON retrospectives(created_by);
CREATE INDEX idx_retrospective_items_retrospective_id ON retrospective_items(retrospective_id);
CREATE INDEX idx_retrospective_items_category ON retrospective_items(category);
CREATE INDEX idx_retrospective_votes_item_id ON retrospective_votes(item_id);
CREATE INDEX idx_retrospective_votes_user_id ON retrospective_votes(user_id);
CREATE INDEX idx_action_items_retrospective_id ON action_items(retrospective_id);
CREATE INDEX idx_action_items_assigned_to ON action_items(assigned_to);
CREATE INDEX idx_action_items_status ON action_items(status);
