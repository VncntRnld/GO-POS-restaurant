package repositories

import (
	"context"
	"database/sql"
	"errors"
	"pos-restaurant/models"
)

type OutletRepository struct {
	db *sql.DB
}

func NewOutletRepository(db *sql.DB) *OutletRepository {
	return &OutletRepository{db: db}
}

func (r *OutletRepository) Create(ctx context.Context, outlet *models.Outlet) (int, error) {
	query := `INSERT INTO outlets (name, location, service_charge_percentage, tax_percentage, is_active) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query,
		outlet.Name, outlet.Location, outlet.ServiceChargePercent, outlet.TaxPercentage, outlet.IsActive,
	).Scan(&outlet.ID)
	return outlet.ID, err
}

func (r *OutletRepository) List(ctx context.Context) ([]*models.Outlet, error) {
	query := `SELECT id, name, location, service_charge_percentage, tax_percentage, is_active 
	          FROM outlets WHERE deleted_at IS NULL ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var outlets []*models.Outlet
	for rows.Next() {
		var outlet models.Outlet
		err := rows.Scan(&outlet.ID, &outlet.Name, &outlet.Location,
			&outlet.ServiceChargePercent, &outlet.TaxPercentage,
			&outlet.IsActive)
		if err != nil {
			return nil, err
		}
		outlets = append(outlets, &outlet)
	}
	return outlets, nil
}

func (r *OutletRepository) GetByID(ctx context.Context, id int) (*models.Outlet, error) {
	query := `SELECT id, name, location, service_charge_percentage, tax_percentage, is_active 
	          FROM outlets WHERE id = $1 AND deleted_at IS NULL`
	row := r.db.QueryRowContext(ctx, query, id)

	var outlet models.Outlet
	err := row.Scan(&outlet.ID, &outlet.Name, &outlet.Location,
		&outlet.ServiceChargePercent, &outlet.TaxPercentage,
		&outlet.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &outlet, nil
}

func (r *OutletRepository) Update(ctx context.Context, outlet *models.Outlet) error {
	query := `UPDATE outlets SET name=$1, location=$2, service_charge_percentage=$3, 
	          tax_percentage=$4, is_active=$5, updated_at=NOW() WHERE id=$6 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query,
		outlet.Name, outlet.Location, outlet.ServiceChargePercent, outlet.TaxPercentage, outlet.IsActive, outlet.ID)
	return err
}

func (r *OutletRepository) SoftDelete(ctx context.Context, id int) error {
	query := `UPDATE outlets SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
