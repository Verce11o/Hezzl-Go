package service

import (
	"context"
	"github.com/Verce11o/Hezzl-Go/api"
	"github.com/Verce11o/Hezzl-Go/internal/models"
	"github.com/Verce11o/Hezzl-Go/internal/product"
	"go.uber.org/zap"
)

type Service struct {
	log  *zap.SugaredLogger
	repo product.Repository
}

func NewService(log *zap.SugaredLogger, repo product.Repository) *Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) CreateProduct(ctx context.Context, projectID int, input api.CreateProductJSONBody) (models.Product, error) {
	prod, err := s.repo.CreateProduct(ctx, projectID, input)

	if err != nil {
		s.log.Errorf("cannot create product: %v", err)
		return models.Product{}, err
	}

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
	_, err := s.GetProductByID(ctx, productID)

	if err != nil {
		return err
	}

	err = s.repo.DeleteProduct(ctx, productID, projectID)

	if err != nil {
		s.log.Errorf("cannot delete product: %v", err)
		return err
	}

	return nil
}

func (s *Service) GetProductsList(ctx context.Context, limit, offset int) (models.ProductList, error) {
	products, err := s.repo.GetProductsList(ctx, limit, offset)
	if err != nil {
		s.log.Errorf("cannot get product list: %v", err)
		return models.ProductList{}, err
	}

	return products, nil
}

func (s *Service) UpdateProductPriority(ctx context.Context, productID, projectID, priority int) ([]models.Priority, error) {
	_, err := s.GetProductByID(ctx, productID)

	if err != nil {
		return nil, err
	}

	products, err := s.repo.UpdateProductPriority(ctx, productID, projectID, priority)
	if err != nil {
		s.log.Errorf("cannot update product priority: %v", err)
		return nil, err
	}

	return products, nil
}
