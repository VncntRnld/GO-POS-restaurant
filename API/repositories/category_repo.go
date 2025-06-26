package repositories

import (
	"context"
	"database/sql"
	"pos-restaurant/models"
)

type MenuCategoryRepository struct {
	db *sql.DB
}

func NewMenuCategoryRepository(db *sql.DB) *MenuCategoryRepository {
	return &MenuCategoryRepository{db: db}
}

func (r *MenuCategoryRepository) Create(ctx context.Context, item *models.MenuCategory) (int, error) {

	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO menu_categories (
			name
		) VALUES ($1)
		RETURNING id
	`,
		item.Name,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *MenuCategoryRepository) List(ctx context.Context) ([]*models.MenuCategory, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM menu_categories WHERE deleted_at IS NULL ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.MenuCategory
	for rows.Next() {
		var c models.MenuCategory
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}

	return categories, nil
}

func (r *MenuCategoryRepository) SoftDelete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE menu_categories
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)

	return err
}
