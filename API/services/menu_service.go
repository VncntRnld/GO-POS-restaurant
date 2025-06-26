package services

import (
	"context"
	"pos-restaurant/models"
	"pos-restaurant/repositories"
)

type MenuService struct {
	repo *repositories.MenuItemRepository
}

func NewMenuService(repo *repositories.MenuItemRepository) *MenuService {
	return &MenuService{repo: repo}
}

func (s *MenuService) CreateMenuItem(ctx context.Context, item *models.MenuItem) (int, error) {
	if item.Cost > item.Price {
		item.IsActive = false // Nonaktifkan jika tidak menguntungkan : Contoh
	}

	return s.repo.Create(ctx, item)
}

func (s *MenuService) ListMenuItemsById(ctx context.Context, id int) (*models.MenuItem, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *MenuService) ListMenuItems(ctx context.Context) ([]*models.MenuItem, error) {
	return s.repo.List(ctx)
}

func (s *MenuService) ListActiveMenuItems(ctx context.Context) ([]*models.MenuItem, error) {
	return s.repo.ListActive(ctx)
}

func (s *MenuService) SearchMenuItems(ctx context.Context, keyword string) ([]*models.MenuItem, error) {
	return s.repo.SearchName(ctx, keyword)
}

func (s *MenuService) UpdateMenuItem(ctx context.Context, item *models.MenuItem) error {
	return s.repo.Update(ctx, item)
}

func (s *MenuService) DeleteMenuItem(ctx context.Context, id int) error {
	return s.repo.SoftDelete(ctx, id)
}
