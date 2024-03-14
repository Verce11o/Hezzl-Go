-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goods (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    priority INT DEFAULT 0,
    removed BOOL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE INDEX idx_product_id ON goods(id);
CREATE INDEX idx_product_project_id ON goods(project_id);
CREATE INDEX idx_name ON goods(name);

CREATE OR REPLACE FUNCTION product_priority_max() RETURNS int
AS 'SELECT COALESCE(max(priority),0)+1 FROM goods'
    LANGUAGE 'sql';

alter table goods alter column priority set default product_priority_max();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE goods;
-- +goose StatementEnd
