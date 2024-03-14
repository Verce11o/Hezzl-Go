package clickhouse

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/Verce11o/Hezzl-Go/internal/config"
	"log"
)

func NewClickHouseClient(ctx context.Context, cfg *config.Config) driver.Conn {

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.ClickHouse.Addr},
		Auth: clickhouse.Auth{
			Database: cfg.ClickHouse.Database,
			Username: cfg.ClickHouse.Username,
			Password: cfg.ClickHouse.Password,
		},
	})

	if err != nil {
		log.Fatalf("cannot connect to clickhouse: %v", err)
	}

	if err = conn.Ping(ctx); err != nil {
		log.Fatalf("cannot connect to clickhouse: %v", err)
	}

	return conn

}
