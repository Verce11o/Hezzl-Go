-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goods (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    priority INT DEFAULT 0,
    removed BOOL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_product_id ON goods(id);
CREATE INDEX idx_product_project_id ON goods(project_id);
CREATE INDEX idx_name ON goods(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE goods;
-- +goose StatementEnd
