package service

import (
	"context"
	"github.com/Verce11o/Hezzl-Go/api"
	"github.com/Verce11o/Hezzl-Go/internal/models"
	"github.com/Verce11o/Hezzl-Go/internal/product"
	"github.com/Verce11o/Hezzl-Go/internal/product/handler/nats"
	"go.uber.org/zap"
)

type Service struct {
	log       *zap.SugaredLogger
	repo      product.Repository
	redis     product.RedisRepository
	publisher *nats.Publisher
}

func NewService(log *zap.SugaredLogger, repo product.Repository, redis product.RedisRepository, publisher *nats.Publisher) *Service {
	return &Service{log: log, repo: repo, redis: redis, publisher: publisher}
}

func (s *Service) CreateProduct(ctx context.Context, projectID int, input api.CreateProductJSONBody) (models.Product, error) {
	prod, err := s.repo.CreateProduct(ctx, projectID, input)

	if err != nil {
		s.log.Errorf("cannot create product: %v", err)
		return models.Product{}, err
	}

	if err = s.redis.DeleteProductList(ctx); err != nil {
		s.log.Errorf("cannot clear redis: %v", err)
	}

	s.publisher.Publish(ctx, prod, "NEW")

	return prod, nil
}

func (s *Service) UpdateProduct(ctx context.Context, productID, projectID int, input api.UpdateProductJSONBody) (models.Product, error) {
	_, err := s.GetProductByID(ctx, productID)

	if err != nil {
		return models.Product{}, err
	}

	prod, err := s.repo.UpdateProduct(ctx, productID, projectID, input)

	if err != nil {
		s.log.Errorf("cannot update product: %v", err)
		return models.Product{}, err
	}

	if err = s.redis.DeleteProductList(ctx); err != nil {
		s.log.Errorf("cannot clear redis: %v", err)
	}

	s.publisher.Publish(ctx, prod, "UPDATE")

	return prod, nil
}

func (s *Service) GetProductByID(ctx context.Context, productID int) (models.Product, error) {
	prod, err := s.repo.GetProductByID(ctx, productID)

	if err != nil {
		s.log.Errorf("cannot get product by id: %v", err)
		return models.Product{}, err
	}

	return prod, nil
}

func (s *Service) DeleteProduct(ctx context.Context, productID, projectID int) error {
	prod, err := s.GetProductByID(ctx, productID)

	if err != nil {
		return err
	}

	err = s.repo.DeleteProduct(ctx, productID, projectID)

	if err != nil {
		s.log.Errorf("cannot delete product: %v", err)
		return err
	}

	if err = s.redis.DeleteProductList(ctx); err != nil {
		s.log.Errorf("cannot clear redis: %v", err)
	}

	s.publisher.Publish(ctx, prod, "DELETE")
	return nil
}

func (s *Service) GetProductsList(ctx context.Context, limit, offset int) (models.ProductList, error) {
	cachedProducts, err := s.redis.GetProductList(ctx, limit, offset)

	if err != nil {
		s.log.Infof("cannot get product list in redis: %v", err)
	}

	if cachedProducts != nil {
		return *cachedProducts, nil
	}

	products, err := s.repo.GetProductsList(ctx, limit, offset)
	if err != nil {
		s.log.Errorf("cannot get product list: %v", err)
		return models.ProductList{}, err
	}

	if err = s.redis.SetByIDCtx(ctx, limit, offset, products); err != nil {
		s.log.Errorf("cannot set product list in redis: %v", err)
	}

	return products, nil
}

func (s *Service) UpdateProductPriority(ctx context.Context, productID, projectID, priority int) ([]models.Priority, error) {
	prod, err := s.GetProductByID(ctx, productID)

	if err != nil {
		return nil, err
	}

	products, err := s.repo.UpdateProductPriority(ctx, productID, projectID, priority)
	if err != nil {
		s.log.Errorf("cannot update product priority: %v", err)
		return nil, err
	}

	s.publisher.Publish(ctx, prod, "UPDATE")

	return products, nil
}
