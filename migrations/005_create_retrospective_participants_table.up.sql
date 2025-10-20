CREATE TABLE retrospective_participants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    retrospective_id UUID NOT NULL REFERENCES retrospectives(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(retrospective_id, user_id)
);

CREATE INDEX idx_retrospective_participants_retrospective_id ON retrospective_participants(retrospective_id);
CREATE INDEX idx_retrospective_participants_user_id ON retrospective_participants(user_id);
