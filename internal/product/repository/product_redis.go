package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Verce11o/Hezzl-Go/internal/models"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	productTTL = 60
)

type ProductRedis struct {
	client *redis.Client
}

func NewProductsRedis(client *redis.Client) *ProductRedis {
	return &ProductRedis{client: client}
}

func (r *ProductRedis) GetProductList(ctx context.Context, limit, offset int) (*models.ProductList, error) {

	productListBytes, err := r.client.Get(ctx, r.createKey(limit, offset)).Bytes()

	if err != nil {
		return nil, err
	}

	var productList models.ProductList

	if err = json.Unmarshal(productListBytes, &productList); err != nil {
		return nil, err
	}

	return &productList, nil
}

func (r *ProductRedis) SetByIDCtx(ctx context.Context, limit, offset int, products models.ProductList) error {

	productListBytes, err := json.Marshal(products)

	if err != nil {
		return err
	}

	return r.client.Set(ctx, r.createKey(limit, offset), productListBytes, time.Second*time.Duration(productTTL)).Err()
}

func (r *ProductRedis) DeleteProductList(ctx context.Context) error {

	return r.client.FlushDB(ctx).Err()
}

func (r *ProductRedis) createKey(limit, offset int) string {
	return fmt.Sprintf("product:%d,%d", limit, offset)
}
