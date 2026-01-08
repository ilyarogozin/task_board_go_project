CREATE TABLE boards (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL
);

CREATE TABLE columns (
    id UUID PRIMARY KEY,
    board_id UUID REFERENCES boards(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    position INT NOT NULL
);

CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    column_id UUID REFERENCES columns(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    assignee_id UUID,
    position INT NOT NULL
);

CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    processed BOOLEAN DEFAULT FALSE
);