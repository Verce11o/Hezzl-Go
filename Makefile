postgres-up:
	goose -dir ./migrations/postgres postgres "host=localhost password=password user=postgres port=5432 dbname=postgres sslmode=disable" up
postgres-down:
	goose -dir ./migrations/postgres postgres "host=localhost password=password user=postgres port=5432 dbname=postgres sslmode=disable" down

clickhouse-up:
	goose -dir ./migrations/clickhouse clickhouse "tcp://127.0.0.1:19000" up
clickhouse-down:
	goose -dir ./migrations/clickhouse clickhouse "tcp://127.0.0.1:19000" down