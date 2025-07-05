package services

import (
	"context"
	"pos-restaurant/models"
	"pos-restaurant/repositories"
)

type TableTransferService struct {
	repo *repositories.TableTransferRepository
}

func NewTableTransferService(repo *repositories.TableTransferRepository) *TableTransferService {
	return &TableTransferService{repo}
}

func (s *TableTransferService) Create(ctx context.Context, t *models.TableTransfer) (int, error) {
	return s.repo.Create(ctx, t)
}

func (s *TableTransferService) List(ctx context.Context) ([]models.TableTransfer, error) {
	return s.repo.List(ctx)
}

func (s *TableTransferService) GetByID(ctx context.Context, id int) (*models.TableTransfer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TableTransferService) Update(ctx context.Context, t *models.TableTransfer) error {
	return s.repo.Update(ctx, t)
}

func (s *TableTransferService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
