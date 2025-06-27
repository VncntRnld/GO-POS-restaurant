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
		INSERT INTO ingredients (name, qty, unit, is_allergen, is_active, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		ing.Name, ing.Qty, ing.Unit, ing.IsAllergen, ing.IsActive, ing.Description,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *IngredientRepository) List(ctx context.Context) ([]*models.Ingredient, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, qty, unit, is_allergen, is_active, description FROM ingredients
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

		err := rows.Scan(&ing.ID, &ing.Name, &ing.Qty, &ing.Unit, &ing.IsAllergen, &ing.IsActive, &desc)
		if err != nil {
			return nil, err
		}
		ing.Description = desc
		ingredients = append(ingredients, &ing)
	}
	return ingredients, nil
}

func (r *IngredientRepository) GetByID(ctx context.Context, id int) (*models.Ingredient, error) {
	query := `
		SELECT id, name, qty, unit, is_allergen, is_active, description
		FROM ingredients
		WHERE id = $1 AND deleted_at IS NULL
	`

	row := r.db.QueryRowContext(ctx, query, id)
	var ingredient models.Ingredient

	err := row.Scan(
		&ingredient.ID,
		&ingredient.Name,
		&ingredient.Qty,
		&ingredient.Unit,
		&ingredient.IsAllergen,
		&ingredient.IsActive,
		&ingredient.Description,
	)

	if err != nil {
		return nil, err
	}

	return &ingredient, nil
}

func (r *IngredientRepository) Update(ctx context.Context, ing *models.Ingredient) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE ingredients SET 
			name = $1,
			qty = $2,
			unit = $3,
			is_allergen = $4,
			is_active = $5,
			description = $6
		WHERE id = $7 AND deleted_at IS NULL
	`, ing.Name, ing.Qty, ing.Unit, ing.IsAllergen, ing.IsActive, ing.Description, ing.ID)

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
