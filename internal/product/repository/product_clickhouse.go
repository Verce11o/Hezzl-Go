package repository

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/Verce11o/Hezzl-Go/internal/models"
)

type ProductClickHouse struct {
	click driver.Conn
}

func NewProductClickHouse(click driver.Conn) *ProductClickHouse {
	return &ProductClickHouse{click: click}
}

func (p *ProductClickHouse) UploadEvent(ctx context.Context, data []models.Product) error {
	batch, err := p.click.PrepareBatch(ctx, "INSERT INTO events")

	if err != nil {
		return err
	}

	for _, val := range data {
		err = batch.Append(
			int32(val.ID),
			int32(val.ProjectID),
			val.Name,
			val.Description,
			int32(val.Priority),
			val.Removed,
			val.CreatedAt,
		)
		if err != nil {
			return err
		}
	}
	return batch.Send()
}
