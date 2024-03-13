migrate-up:
	goose -dir ./migrations postgres "host=localhost password=vercello user=postgres port=5432 dbname=hezzl sslmode=disable" up
migrate-down:
	goose -dir ./migrations postgres "host=localhost password=vercello user=postgres port=5432 dbname=hezzl sslmode=disable" down