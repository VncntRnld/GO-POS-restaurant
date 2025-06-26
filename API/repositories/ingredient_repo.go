package repositories

import (
	"context"
	"database/sql"
	"pos-restaurant/models"
)

type IngredientRepository struct {
	db *sql.DB
}

func NewIngredientRepository(db *sql.DB) *IngredientRepository {
	return &IngredientRepository{db: db}
}

func (r *IngredientRepository) Create(ctx context.Context, ing *models.Ingredient) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO ingredients (name, qty, is_allergen, is_active, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		ing.Name, ing.Qty, ing.IsAllergen, ing.IsActive, ing.Description,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *IngredientRepository) List(ctx context.Context) ([]*models.Ingredient, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, qty, is_allergen, is_active, description FROM ingredients
		WHERE deleted_at IS NULL
		ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []*models.Ingredient
	for rows.Next() {
		var ing models.Ingredient
		var desc sql.NullString

		err := rows.Scan(&ing.ID, &ing.Name, &ing.Qty, &ing.IsAllergen, &ing.IsActive, &desc)
		if err != nil {
			return nil, err
		}
		ing.Description = desc
		ingredients = append(ingredients, &ing)
	}
	return ingredients, nil
}

func (r *IngredientRepository) Update(ctx context.Context, ing *models.Ingredient) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ingredients SET 
			name = $1,
			qty = $2,
			is_allergen = $3,
			is_active = $4,
			description = $5
		WHERE id = $6 AND deleted_at IS NULL
	`, ing.Name, ing.Qty, ing.IsAllergen, ing.IsActive, ing.Description, ing.ID)

	return err
}

func (r *IngredientRepository) SoftDelete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ingredients 
		SET deleted_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL
	`, id)

	return err
}
