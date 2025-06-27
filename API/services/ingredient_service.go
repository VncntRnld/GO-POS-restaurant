package services

import (
	"context"
	"pos-restaurant/models"
	"pos-restaurant/repositories"
)

type IngredientService struct {
	repo *repositories.IngredientRepository
}

func NewIngredientService(repo *repositories.IngredientRepository) *IngredientService {
	return &IngredientService{repo: repo}
}

func (s *IngredientService) CreateIngredient(ctx context.Context, ing *models.Ingredient) (int, error) {
	return s.repo.Create(ctx, ing)
}

func (s *IngredientService) ListIngredients(ctx context.Context) ([]*models.Ingredient, error) {
	return s.repo.List(ctx)
}

func (s *IngredientService) GetIngredientByID(ctx context.Context, id int) (*models.Ingredient, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *IngredientService) UpdateIngredient(ctx context.Context, ing *models.Ingredient) error {
	return s.repo.Update(ctx, ing)
}

func (s *IngredientService) DeleteIngredient(ctx context.Context, id int) error {
	return s.repo.SoftDelete(ctx, id)
}
