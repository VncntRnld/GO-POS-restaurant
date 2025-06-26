package services

import (
	"context"
	"pos-restaurant/models"
	"pos-restaurant/repositories"
)

type MenuCategoryService struct {
	repo *repositories.MenuCategoryRepository
}

func NewMenuCategoryService(repo *repositories.MenuCategoryRepository) *MenuCategoryService {
	return &MenuCategoryService{repo: repo}
}

func (s *MenuCategoryService) CreateCategory(ctx context.Context, category *models.MenuCategory) (int, error) {
	return s.repo.Create(ctx, category)
}

func (s *MenuCategoryService) ListCategories(ctx context.Context) ([]*models.MenuCategory, error) {
	return s.repo.List(ctx)
}

func (s *MenuCategoryService) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.SoftDelete(ctx, id)
}
