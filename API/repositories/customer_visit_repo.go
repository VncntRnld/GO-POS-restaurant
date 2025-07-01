package repositories

import (
	"context"
	"database/sql"
	"pos-restaurant/models"
)

type CustomerVisitRepository struct {
	db *sql.DB
}

func NewCustomerVisitRepository(db *sql.DB) *CustomerVisitRepository {
	return &CustomerVisitRepository{db: db}
}

func (r *CustomerVisitRepository) Create(ctx context.Context, visit *models.CustomerVisit) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO customer_visits (
			customer_id, visit_type, visit_date, room_number, reservation_id,
			outlet_id, total_spent, pax
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id
	`,
		visit.CustomerID, visit.VisitType, visit.VisitDate,
		visit.RoomNumber, visit.ReservationID, visit.OutletID,
		visit.TotalSpent, visit.Pax,
	).Scan(&id)

	return id, err
}

func (r *CustomerVisitRepository) List(ctx context.Context) ([]*models.CustomerVisit, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, customer_id, visit_type, visit_date,
		       room_number, reservation_id, outlet_id, total_spent, pax
		FROM customer_visits
		ORDER BY visit_date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visits []*models.CustomerVisit
	for rows.Next() {
		var visit models.CustomerVisit
		err := rows.Scan(
			&visit.ID, &visit.CustomerID, &visit.VisitType, &visit.VisitDate,
			&visit.RoomNumber, &visit.ReservationID, &visit.OutletID,
			&visit.TotalSpent, &visit.Pax,
		)
		if err != nil {
			return nil, err
		}
		visits = append(visits, &visit)
	}
	return visits, nil
}

func (r *CustomerVisitRepository) GetByID(ctx context.Context, id int) (*models.CustomerVisit, error) {
	var visit models.CustomerVisit
	err := r.db.QueryRowContext(ctx, `
		SELECT id, customer_id, visit_type, visit_date,
		       room_number, reservation_id, outlet_id, total_spent, pax
		FROM customer_visits
		WHERE id = $1
	`, id).Scan(
		&visit.ID, &visit.CustomerID, &visit.VisitType, &visit.VisitDate,
		&visit.RoomNumber, &visit.ReservationID, &visit.OutletID,
		&visit.TotalSpent, &visit.Pax,
	)
	if err != nil {
		return nil, err
	}
	return &visit, nil
}

func (r *CustomerVisitRepository) Update(ctx context.Context, visit *models.CustomerVisit) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE customer_visits SET
			customer_id = $1, visit_type = $2, visit_date = $3,
			room_number = $4, reservation_id = $5, outlet_id = $6,
			total_spent = $7, pax = $8, updated_at = NOW()
		WHERE id = $9
	`,
		visit.CustomerID, visit.VisitType, visit.VisitDate,
		visit.RoomNumber, visit.ReservationID, visit.OutletID,
		visit.TotalSpent, visit.Pax, visit.ID,
	)
	return err
}

func (r *CustomerVisitRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM customer_visits WHERE id = $1`, id)
	return err
}
