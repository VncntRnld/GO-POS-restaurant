package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"pos-restaurant/models"
)

type MenuItemRepository struct {
	db *sql.DB
}

func NewMenuItemRepository(db *sql.DB) *MenuItemRepository {
	return &MenuItemRepository{db: db}
}

func (r *MenuItemRepository) Create(ctx context.Context, item *models.MenuItem) (int, error) {
	tagsJSON, _ := json.Marshal(item.Tags)

	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO menu_items (
			category_id, sku, name, description, price, cost, 
			is_active, preparation_time, tags
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`,
		item.CategoryID,
		item.SKU,
		item.Name,
		item.Description,
		item.Price,
		item.Cost,
		item.IsActive,
		item.PreparationTime,
		tagsJSON,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *MenuItemRepository) Update(ctx context.Context, item *models.MenuItem) error {
	tagsJSON, _ := json.Marshal(item.Tags)

	_, err := r.db.ExecContext(ctx, `
		UPDATE menu_items SET
			category_id = $1,
			sku = $2,
			name = $3,
			description = $4,
			price = $5,
			cost = $6,
			is_active = $7,
			preparation_time = $8,
			tags = $9,
			updated_at = NOW()
		WHERE id = $10 AND deleted_at IS NULL
	`,
		item.CategoryID,
		item.SKU,
		item.Name,
		item.Description,
		item.Price,
		item.Cost,
		item.IsActive,
		item.PreparationTime,
		tagsJSON,
		item.ID,
	)

	return err
}

func (r *MenuItemRepository) GetByID(ctx context.Context, id int) (*models.MenuItem, error) {
	item := &models.MenuItem{}
	var tagsJSON []byte
	var description sql.NullString
	var prepTime sql.NullInt64

	err := r.db.QueryRowContext(ctx, `
		SELECT
			id, category_id, sku, name, description, price, cost,
			is_active, preparation_time, tags, created_at, updated_at
		FROM menu_items
		WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(
		&item.ID,
		&item.CategoryID,
		&item.SKU,
		&item.Name,
		&description,
		&item.Price,
		&item.Cost,
		&item.IsActive,
		&prepTime,
		&tagsJSON,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("menu item not found")
		}
		return nil, err
	}

	// Handle nullable fields
	item.Description = description
	item.PreparationTime = prepTime

	// Parse JSON tags
	if len(tagsJSON) > 0 {
		json.Unmarshal(tagsJSON, &item.Tags)
	}

	return item, nil
}

func (r *MenuItemRepository) List(ctx context.Context) ([]*models.MenuItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT 
			id, category_id, sku, name, description, price, cost,
			is_active, preparation_time, tags
		FROM menu_items
		WHERE deleted_at IS NULL
		ORDER BY name
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []*models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		var tagsJSON []byte
		var description sql.NullString
		var prepTime sql.NullInt64

		err := rows.Scan(
			&item.ID,
			&item.CategoryID,
			&item.SKU,
			&item.Name,
			&description,
			&item.Price,
			&item.Cost,
			&item.IsActive,
			&prepTime,
			&tagsJSON,
		)
		if err != nil {
			return nil, err
		}

		item.Description = description
		item.PreparationTime = prepTime
		json.Unmarshal(tagsJSON, &item.Tags)

		items = append(items, &item)
	}

	return items, nil
}

func (r *MenuItemRepository) ListActive(ctx context.Context) ([]*models.MenuItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT 
			id, category_id, sku, name, description, price, cost,
			is_active, preparation_time, tags
		FROM menu_items
		WHERE is_active = TRUE AND deleted_at IS NULL
		ORDER BY name
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []*models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		var tagsJSON []byte
		var description sql.NullString
		var prepTime sql.NullInt64

		err := rows.Scan(
			&item.ID,
			&item.CategoryID,
			&item.SKU,
			&item.Name,
			&description,
			&item.Price,
			&item.Cost,
			&item.IsActive,
			&prepTime,
			&tagsJSON,
		)
		if err != nil {
			return nil, err
		}

		item.Description = description
		item.PreparationTime = prepTime
		json.Unmarshal(tagsJSON, &item.Tags)

		items = append(items, &item)
	}

	return items, nil
}

func (r *MenuItemRepository) SearchName(ctx context.Context, name string) ([]*models.MenuItem, error) {
	query := `
		SELECT id, category_id, sku, name, description, price, cost,
			   is_active, preparation_time, tags
		FROM menu_items
		WHERE LOWER(name) LIKE LOWER($1) AND deleted_at IS NULL
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		var tagsJSON []byte
		var description sql.NullString
		var prepTime sql.NullInt64

		err := rows.Scan(
			&item.ID,
			&item.CategoryID,
			&item.SKU,
			&item.Name,
			&description,
			&item.Price,
			&item.Cost,
			&item.IsActive,
			&prepTime,
			&tagsJSON,
		)
		if err != nil {
			return nil, err
		}

		item.Description = description
		item.PreparationTime = prepTime
		json.Unmarshal(tagsJSON, &item.Tags)

		items = append(items, &item)
	}

	return items, nil
}

func (r *MenuItemRepository) SoftDelete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE menu_items
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)

	if err != nil {
		return err
	}

	return nil
}
