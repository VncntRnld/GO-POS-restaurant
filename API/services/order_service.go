package services

import (
	"context"
	"pos-restaurant/models"
	"pos-restaurant/repositories"

	"github.com/google/uuid"
)

type OrderService struct {
	repo *repositories.OrderRepository
}

func NewOrderService(repo *repositories.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(ctx context.Context, req *models.OrderRequest) (int, error) {
	req.OrderNumber = uuid.NewString()
	return s.repo.Create(ctx, req)
}

func (s *OrderService) List(ctx context.Context) ([]*models.OrderRequest, error) {
	return s.repo.List(ctx)
}

func (s *OrderService) GetByID(ctx context.Context, id int) (*models.OrderRequest, error) {
	return s.repo.GetByID(ctx, id)
}

// Belum
func (s *OrderService) Update(ctx context.Context, order *models.Order) error {
	return s.repo.Update(ctx, order)
}

func (s *OrderService) SoftDelete(ctx context.Context, id int) error {
	return s.repo.SoftDelete(ctx, id)
}
