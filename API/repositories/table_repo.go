package repositories

import (
	"context"
	"database/sql"
	"pos-restaurant/models"
)

type TableRepository struct {
	db *sql.DB
}

func NewTableRepository(db *sql.DB) *TableRepository {
	return &TableRepository{db: db}
}

func (r *TableRepository) Create(ctx context.Context, table *models.Table) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO tables (outlet_id, table_number, capacity, location_type, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		table.OutletID, table.TableNumber, table.Capacity, table.LocationType, table.Status,
	).Scan(&id)
	return id, err
}

func (r *TableRepository) List(ctx context.Context) ([]*models.Table, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, outlet_id, table_number, capacity, location_type, status
		FROM tables WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*models.Table
	for rows.Next() {
		var t models.Table
		err := rows.Scan(
			&t.ID, &t.OutletID, &t.TableNumber, &t.Capacity, &t.LocationType,
			&t.Status,
		)
		if err != nil {
			return nil, err
		}
		tables = append(tables, &t)
	}
	return tables, nil
}

func (r *TableRepository) GetByID(ctx context.Context, id int) (*models.Table, error) {
	var t models.Table
	err := r.db.QueryRowContext(ctx, `
		SELECT id, outlet_id, table_number, capacity, location_type, status
		FROM tables WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(
			&t.ID, &t.OutletID, &t.TableNumber, &t.Capacity, &t.LocationType,
			&t.Status,
		)
	return &t, err
}

func (r *TableRepository) Update(ctx context.Context, table *models.Table) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tables SET outlet_id = $1, table_number = $2, capacity = $3,
		                 location_type = $4, status = $5, updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL`,
		table.OutletID, table.TableNumber, table.Capacity,
		table.LocationType, table.Status, table.ID)
	return err
}

func (r *TableRepository) SoftDelete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tables SET deleted_at = NOW() WHERE id = $1`, id)
	return err
}
