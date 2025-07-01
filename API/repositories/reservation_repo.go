package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"pos-restaurant/models"
)

type ReservationRepository struct {
	db *sql.DB
}

func NewReservationRepository(db *sql.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) Create(ctx context.Context, res *models.Reservation) (int, error) {

	// Cek double booking
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM reservations
			WHERE table_id = $1 AND reservation_time = $2
		)
	`, res.TableID, res.ReservationTime).Scan(&exists)

	if err != nil {
		return 0, err
	}

	if exists {
		return 0, fmt.Errorf("meja sudah dipesan pada waktu tersebut")
	}

	var id int
	err = r.db.QueryRowContext(ctx, `
		INSERT INTO reservations (
			customer_id, reservation_time, pax, table_id, status, special_request
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`,
		res.CustomerID,
		res.ReservationTime,
		res.Pax,
		res.TableID,
		res.Status,
		res.SpecialRequest,
	).Scan(&id)

	return id, err
}

func (r *ReservationRepository) List(ctx context.Context, sortBy string) ([]*models.ReservationWithDetails, error) {
	validSorts := map[string]string{
		"reservation_time": "r.reservation_time",
		"customer_name":    "c.name",
		"status":           "r.status",
		"table":            "t.table_number",
	}

	sortColumn, ok := validSorts[sortBy]
	if !ok {
		sortColumn = "r.reservation_time" // default
	}

	query := fmt.Sprintf(`
		SELECT 
			r.id, r.customer_id, c.name, r.reservation_time, r.pax, 
			t.table_number, r.status, r.special_request
		FROM reservations r
		LEFT JOIN customers c ON r.customer_id = c.cust_id
		LEFT JOIN "tables" t ON r.table_id = t.id
		ORDER BY %s
	`, sortColumn)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*models.ReservationWithDetails
	for rows.Next() {
		var res models.ReservationWithDetails
		var specialReq sql.NullString
		err := rows.Scan(
			&res.ID, &res.CustomerName, &res.CustomerName,
			&res.ReservationTime, &res.Pax,
			&res.TableNumber, &res.Status, &specialReq,
		)
		if err != nil {
			return nil, err
		}
		res.SpecialRequest = specialReq
		reservations = append(reservations, &res)
	}

	return reservations, nil
}

func (r *ReservationRepository) GetByID(ctx context.Context, id int) (*models.Reservation, error) {
	var res models.Reservation
	err := r.db.QueryRowContext(ctx, `
		SELECT id, customer_id, reservation_time, pax, table_id, status, special_request
		FROM reservations WHERE id = $1
	`, id).Scan(
		&res.ID,
		&res.CustomerID,
		&res.ReservationTime,
		&res.Pax,
		&res.TableID,
		&res.Status,
		&res.SpecialRequest,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *ReservationRepository) Update(ctx context.Context, res *models.Reservation) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE reservations SET
			customer_id = $1, reservation_time = $2, pax = $3, table_id = $4,
			status = $5, special_request = $6, updated_at = NOW()
		WHERE id = $7
	`,
		res.CustomerID,
		res.ReservationTime,
		res.Pax,
		res.TableID,
		res.Status,
		res.SpecialRequest,
		res.ID,
	)
	return err
}

func (r *ReservationRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM reservations WHERE id = $1`, id)
	return err
}
