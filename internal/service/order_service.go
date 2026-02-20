package service

import (
	"context"

	"github.com/Dawwami/go-order-api/internal/model"
	"github.com/Dawwami/go-order-api/internal/repository"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(ctx context.Context, order *model.Order) error {
	return s.repo.Create(ctx, order)
}

func (s *OrderService) GetAll(ctx context.Context) ([]model.Order, error) {
	return s.repo.FindAll(ctx)
}

func (s *OrderService) GetByID(ctx context.Context, id uint) (*model.Order, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *OrderService) GetByUserID(ctx context.Context, userID uint) ([]model.Order, error) {
	return s.repo.FindByUserID(ctx, userID)
}
