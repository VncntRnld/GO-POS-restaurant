package repositories

import (
	"context"
	"database/sql"
	"pos-restaurant/models"
)

type StaffRepository struct {
	db *sql.DB
}

func NewStaffRepository(db *sql.DB) *StaffRepository {
	return &StaffRepository{db: db}
}

func (r *StaffRepository) Create(ctx context.Context, staff *models.Staff) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO staff (name, role, pin_code, is_active)
		VALUES ($1, $2, $3, $4) RETURNING id
	`, staff.Name, staff.Role, staff.PinCode, staff.IsActive).Scan(&id)
	return id, err
}

func (r *StaffRepository) List(ctx context.Context) ([]*models.Staff, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, role, pin_code, is_active
		FROM staff WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Staff
	for rows.Next() {
		var s models.Staff
		err := rows.Scan(&s.ID, &s.Name, &s.Role, &s.PinCode, &s.IsActive)
		if err != nil {
			return nil, err
		}
		result = append(result, &s)
	}
	return result, nil
}

func (r *StaffRepository) GetByID(ctx context.Context, id int) (*models.Staff, error) {
	var s models.Staff
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, role, pin_code, is_active
		FROM staff WHERE id = $1 AND deleted_at IS NULL
	`, id).Scan(&s.ID, &s.Name, &s.Role, &s.PinCode, &s.IsActive)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StaffRepository) Update(ctx context.Context, s *models.Staff) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE staff SET name = $1, role = $2, pin_code = $3, is_active = $4, updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
	`, s.Name, s.Role, s.PinCode, s.IsActive, s.ID)
	return err
}

func (r *StaffRepository) SoftDelete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE staff SET deleted_at = NOW() WHERE id = $1
	`, id)
	return err
}
