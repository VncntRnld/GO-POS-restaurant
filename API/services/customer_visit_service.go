package services

import (
	"context"
	"pos-restaurant/models"
	"pos-restaurant/repositories"
)

type CustomerVisitService struct {
	repo *repositories.CustomerVisitRepository
}

func NewCustomerVisitService(repo *repositories.CustomerVisitRepository) *CustomerVisitService {
	return &CustomerVisitService{repo: repo}
}

func (s *CustomerVisitService) Create(ctx context.Context, visit *models.CustomerVisit) (int, error) {
	return s.repo.Create(ctx, visit)
}

func (s *CustomerVisitService) List(ctx context.Context) ([]*models.CustomerVisit, error) {
	return s.repo.List(ctx)
}

func (s *CustomerVisitService) GetByID(ctx context.Context, id int) (*models.CustomerVisit, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CustomerVisitService) Update(ctx context.Context, visit *models.CustomerVisit) error {
	return s.repo.Update(ctx, visit)
}

func (s *CustomerVisitService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
