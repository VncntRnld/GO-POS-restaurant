package services

import (
	"context"
	"pos-restaurant/models"
	"pos-restaurant/repositories"
)

type TableService struct {
	repo *repositories.TableRepository
}

func NewTableService(repo *repositories.TableRepository) *TableService {
	return &TableService{repo: repo}
}

func (s *TableService) CreateTable(ctx context.Context, table *models.Table) (int, error) {
	return s.repo.Create(ctx, table)
}

func (s *TableService) ListTables(ctx context.Context) ([]*models.Table, error) {
	return s.repo.List(ctx)
}

func (s *TableService) GetTableByID(ctx context.Context, id int) (*models.Table, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TableService) UpdateTable(ctx context.Context, table *models.Table) error {
	return s.repo.Update(ctx, table)
}

func (s *TableService) SoftDeleteTable(ctx context.Context, id int) error {
	return s.repo.SoftDelete(ctx, id)
}
