package services

import (
	"context"
	"pos-restaurant/models"
	"pos-restaurant/repositories"
)

type OutletService struct {
	repo *repositories.OutletRepository
}

func NewOutletService(repo *repositories.OutletRepository) *OutletService {
	return &OutletService{repo: repo}
}

func (s *OutletService) CreateOutlet(ctx context.Context, outlet *models.Outlet) (int, error) {
	return s.repo.Create(ctx, outlet)
}

func (s *OutletService) ListOutlets(ctx context.Context) ([]*models.Outlet, error) {
	return s.repo.List(ctx)
}

func (s *OutletService) GetOutletByID(ctx context.Context, id int) (*models.Outlet, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OutletService) UpdateOutlet(ctx context.Context, outlet *models.Outlet) error {
	return s.repo.Update(ctx, outlet)
}

func (s *OutletService) SoftDeleteOutlet(ctx context.Context, id int) error {
	return s.repo.SoftDelete(ctx, id)
}
