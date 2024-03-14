package repository

import (
	"context"
	"github.com/Verce11o/Hezzl-Go/api"
	"github.com/Verce11o/Hezzl-Go/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (p *ProductRepository) CreateProduct(ctx context.Context, projectID int, input api.CreateProductJSONBody) (models.Product, error) {

	q := "INSERT INTO goods(project_id, name) VALUES ($1, $2) RETURNING id, project_id, name, description, priority, removed, created_at"
	row, err := p.db.Query(ctx, q, projectID, input.Name)

	if err != nil {
		return models.Product{}, err
	}

	prod, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Product])

	if err != nil {
		return models.Product{}, err
	}

	return prod, nil
}

func (p *ProductRepository) UpdateProduct(ctx context.Context, productID, projectID int, input api.UpdateProductJSONBody) (models.Product, error) {

	tx, err := p.db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		return models.Product{}, err
	}

	q := "UPDATE goods SET name = $1, description = COALESCE(NULLIF($2, ''), description) " +
		"WHERE id = $3 AND project_id = $4 RETURNING id, project_id, name, description, priority, removed, created_at"
	row, err := tx.Query(ctx, q, input.Name, input.Description, productID, projectID)

	if err != nil {
		return models.Product{}, err
	}

	prod, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Product])

	if err != nil {
		return models.Product{}, err
	}

	return prod, tx.Commit(ctx)
}

func (p *ProductRepository) DeleteProduct(ctx context.Context, productID, projectID int) error {
	tx, err := p.db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		return err
	}

	q := "UPDATE goods SET removed = TRUE WHERE id = $1 AND project_id = $2"

	_, err = tx.Exec(ctx, q, productID, projectID)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (p *ProductRepository) GetProductByID(ctx context.Context, productID int) (models.Product, error) {
	q := "SELECT id, project_id, name, description, priority, removed, created_at FROM goods WHERE id = $1 AND removed = false"

	row, err := p.db.Query(ctx, q, productID)
	if err != nil {
		return models.Product{}, err
	}

	prod, err := pgx.CollectOneRow(row, pgx.RowToStructByName[models.Product])

	if err != nil {
		return models.Product{}, err
	}

	return prod, nil

}

func (p *ProductRepository) GetProductsList(ctx context.Context, limit, offset int) (models.ProductList, error) {
	goodsQuery := "SELECT id, project_id, name, description, priority, removed, created_at FROM goods ORDER BY id OFFSET COALESCE(NULLIF($1, 0), 1) LIMIT COALESCE(NULLIF($2, 0), 10)"
	rows, err := p.db.Query(ctx, goodsQuery, offset, limit)

	if err != nil {
		return models.ProductList{}, err
	}

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Product])
	if err != nil {
		return models.ProductList{}, err
	}

	meta := models.Meta{
		Total:   len(products),
		Removed: countRemovedGoods(products),
		Limit:   limit,
		Offset:  offset,
	}

	return models.ProductList{Goods: products, Meta: meta}, nil

}

func (p *ProductRepository) UpdateProductPriority(ctx context.Context, productID, projectID, priority int) ([]models.Priority, error) {
	selectRemainGoodsQuery := "SELECT id FROM goods WHERE id >= $1 AND project_id = $2 ORDER BY id"

	rows, err := p.db.Query(ctx, selectRemainGoodsQuery, productID, projectID)
	if err != nil {
		return nil, err
	}

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	maxPriority := priority

	priorities := make([]models.Priority, 0)

	for rows.Next() {
		var rowProductID int

		if err = rows.Scan(&rowProductID); err != nil {
			return nil, err
		}

		_, err = tx.Exec(ctx, "UPDATE goods SET priority = $1 WHERE id = $2 AND project_id = $3", maxPriority, rowProductID, projectID)

		if err != nil {
			return nil, err
		}

		priorities = append(priorities, models.Priority{ID: rowProductID, Priority: maxPriority})
		maxPriority++
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return priorities, nil
}

func countRemovedGoods(goods []models.Product) int {
	count := 0
	for _, prod := range goods {
		if prod.Removed {
			count++
		}
	}
	return count
}
