package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"pos-restaurant/models"
)

type MenuIngredientRepository struct {
	db *sql.DB
}

func NewMenuIngredientRepository(db *sql.DB) *MenuIngredientRepository {
	return &MenuIngredientRepository{db: db}
}

func (r *MenuIngredientRepository) Create(ctx context.Context, mi *models.MenuIngredient) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO menu_ingredients (
			menu_item_id, ingredient_id, qty, is_removable, is_default
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		mi.MenuItemID, mi.IngredientID, mi.Qty, mi.IsRemovable, mi.IsDefault,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *MenuIngredientRepository) ListByMenuItem(ctx context.Context, menuItemID int) ([]*models.MenuIngredient, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, menu_item_id, ingredient_id, qty, is_removable, is_default, created_at, updated_at
		FROM menu_ingredients
		WHERE menu_item_id = $1
	`, menuItemID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.MenuIngredient
	for rows.Next() {
		var mi models.MenuIngredient
		err := rows.Scan(
			&mi.ID, &mi.MenuItemID, &mi.IngredientID,
			&mi.Qty, &mi.IsRemovable, &mi.IsDefault,
			&mi.CreatedAt, &mi.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &mi)
	}

	return result, nil
}

func (r *MenuIngredientRepository) Update(ctx context.Context, m *models.MenuIngredient) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE menu_ingredients SET 
			menu_item_id = $1,
			ingredient_id = $2,
			qty = $3,
			is_removable = $4,
			is_default = $5,
			updated_at = NOW()
		WHERE id = $6
	`, m.MenuItemID, m.IngredientID, m.Qty, m.IsRemovable, m.IsDefault, m.ID)

	return err
}

func (r *MenuIngredientRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM menu_ingredients WHERE id = $1`, id)
	return err
}

func (r *MenuItemRepository) GetMenuWithIngredients(ctx context.Context, id int) (*models.MenuItemWithIngredients, error) {
	query := `
		SELECT 
			mi.id, mi.name, mi.sku, mi.description, mi.price, mi.cost,
			mi.is_active, mi.preparation_time, mi.tags,
			i.id, i.name, mgi.qty, i.unit, i.is_allergen, i.is_active, i.description
		FROM menu_items mi
		LEFT JOIN menu_ingredients mgi ON mi.id = mgi.menu_item_id
		LEFT JOIN ingredients i ON mgi.ingredient_id = i.id
		WHERE mi.id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menu *models.MenuItemWithIngredients
	ingredients := []models.Ingredient{}

	var tagsJSON []byte
	var desc sql.NullString
	var prep sql.NullInt64

	for rows.Next() {
		var ing models.Ingredient
		var ingDesc sql.NullString
		var temp models.MenuItem

		err := rows.Scan(
			&temp.ID, &temp.Name, &temp.SKU, &desc, &temp.Price, &temp.Cost,
			&temp.IsActive, &prep, &tagsJSON,
			&ing.ID, &ing.Name, &ing.Qty, &ing.Unit, &ing.IsAllergen, &ing.IsActive, &ingDesc,
		)
		if err != nil {
			return nil, err
		}

		if menu == nil {
			temp.Description = desc
			temp.PreparationTime = prep
			json.Unmarshal(tagsJSON, &temp.Tags)
			menu = &models.MenuItemWithIngredients{
				MenuItem:    temp,
				Ingredients: []models.Ingredient{},
			}
		}

		if ing.ID != 0 {
			ing.Description = ingDesc
			ingredients = append(ingredients, ing)
		}
	}

	if menu == nil {
		return nil, sql.ErrNoRows
	}

	menu.Ingredients = ingredients
	return menu, nil
}
