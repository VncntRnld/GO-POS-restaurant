package services

import (
	"context"
	"pos-restaurant/models"
	"pos-restaurant/repositories"
)

type MenuIngredientService struct {
	repo *repositories.MenuIngredientRepository
}

func NewMenuIngredientService(repo *repositories.MenuIngredientRepository) *MenuIngredientService {
	return &MenuIngredientService{repo: repo}
}

func (s *MenuIngredientService) CreateMenuIngredient(ctx context.Context, m *models.MenuIngredient) (int, error) {
	return s.repo.Create(ctx, m)
}

func (s *MenuIngredientService) GetIngredientsByMenuItem(ctx context.Context, menuItemID int) ([]*models.MenuIngredient, error) {
	return s.repo.ListByMenuItem(ctx, menuItemID)
}

func (s *MenuIngredientService) UpdateMenuIngredient(ctx context.Context, m *models.MenuIngredient) error {
	return s.repo.Update(ctx, m)
}

func (s *MenuIngredientService) DeleteMenuIngredient(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
