-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE INDEX idx_project_id ON projects(id);
INSERT INTO projects(name) VALUES ('Первая запись');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE projects;
-- +goose StatementEnd
