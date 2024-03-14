-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events
(
    id          int,
    project_id int,
    name        String,
    description String,
    priority    int,
    removed     Bool,
    EventTime   TIMESTAMP,
    INDEX index_id_events id TYPE minmax GRANULARITY 3,
    INDEX index_project_id_events project_id TYPE minmax GRANULARITY 3,
    INDEX index_name_events name TYPE ngrambf_v1(4, 1024, 1, 0 ) GRANULARITY 1
) ENGINE = MergeTree
ORDER BY EventTime;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events
-- +goose StatementEnd
