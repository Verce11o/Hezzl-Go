package product

import (
	"context"
	"github.com/Verce11o/Hezzl-Go/api"
	"github.com/Verce11o/Hezzl-Go/internal/models"
)

type Repository interface {
	CreateProduct(ctx context.Context, projectID int, input api.CreateProductJSONBody) (models.Product, error)
	UpdateProduct(ctx context.Context, productID, projectID int, input api.UpdateProductJSONBody) (models.Product, error)
	GetProductByID(ctx context.Context, productID int) (models.Product, error)
	GetProductsList(ctx context.Context, limit, offset int) (models.ProductList, error)
	UpdateProductPriority(ctx context.Context, productID, projectID, priority int) ([]models.Priority, error)
	DeleteProduct(ctx context.Context, productID, projectID int) error
}

type RedisRepository interface {
	GetProductList(ctx context.Context, limit, offset int) (*models.ProductList, error)
	SetByIDCtx(ctx context.Context, limit, offset int, products models.ProductList) error
	DeleteProductList(ctx context.Context) error
}

type ClickHouseRepository interface {
	UploadEvent(ctx context.Context, data []models.Product) error
}
