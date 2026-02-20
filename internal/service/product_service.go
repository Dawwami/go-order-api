package service

import (
	"context"

	"github.com/Dawwami/go-order-api/internal/model"
	"github.com/Dawwami/go-order-api/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(ctx context.Context) ([]model.Product, error) {
	data, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*model.Product, error) {
	data, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ProductService) Create(ctx context.Context, product *model.Product) error {
	return s.repo.Create(ctx, product)
}

func (s *ProductService) Update(ctx context.Context, product *model.Product) error {
	return s.repo.Update(ctx, product)
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
