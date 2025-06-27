package services

import (
	"context"
	"pos-restaurant/models"
	"pos-restaurant/repositories"
)

type StaffService struct {
	repo *repositories.StaffRepository
}

func NewStaffService(repo *repositories.StaffRepository) *StaffService {
	return &StaffService{repo: repo}
}

func (s *StaffService) CreateStaff(ctx context.Context, staff *models.Staff) (int, error) {
	return s.repo.Create(ctx, staff)
}

func (s *StaffService) ListStaff(ctx context.Context) ([]*models.Staff, error) {
	return s.repo.List(ctx)
}

func (s *StaffService) GetStaffByID(ctx context.Context, id int) (*models.Staff, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *StaffService) UpdateStaff(ctx context.Context, staff *models.Staff) error {
	return s.repo.Update(ctx, staff)
}

func (s *StaffService) SoftDeleteStaff(ctx context.Context, id int) error {
	return s.repo.SoftDelete(ctx, id)
}
