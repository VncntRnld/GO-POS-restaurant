package services

import (
	"context"
	"pos-restaurant/models"
	"pos-restaurant/repositories"
)

type ReservationService struct {
	repo *repositories.ReservationRepository
}

func NewReservationService(repo *repositories.ReservationRepository) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) Create(ctx context.Context, res *models.Reservation) (int, error) {
	return s.repo.Create(ctx, res)
}

func (s *ReservationService) List(ctx context.Context, sortBy string) ([]*models.ReservationWithDetails, error) {
	return s.repo.List(ctx, sortBy)
}

func (s *ReservationService) GetByID(ctx context.Context, id int) (*models.Reservation, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ReservationService) Update(ctx context.Context, res *models.Reservation) error {
	return s.repo.Update(ctx, res)
}

func (s *ReservationService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
