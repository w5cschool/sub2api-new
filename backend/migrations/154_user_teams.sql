-- User teams and leader visibility.

CREATE TABLE IF NOT EXISTS teams (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS teams_name_unique_active
    ON teams (name)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS teams_status
    ON teams (status);

CREATE INDEX IF NOT EXISTS teams_deleted_at
    ON teams (deleted_at);

CREATE TABLE IF NOT EXISTS team_members (
    id         BIGSERIAL PRIMARY KEY,
    team_id    BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role       VARCHAR(20) NOT NULL DEFAULT 'member',
    joined_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT team_members_role_check CHECK (role IN ('member', 'leader'))
);

CREATE UNIQUE INDEX IF NOT EXISTS team_members_user_unique
    ON team_members (user_id);

CREATE INDEX IF NOT EXISTS team_members_team_id
    ON team_members (team_id);

CREATE INDEX IF NOT EXISTS team_members_team_role
    ON team_members (team_id, role);
