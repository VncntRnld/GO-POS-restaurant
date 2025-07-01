// repositories/customer_repository.go
package repositories

import (
	"context"
	"database/sql"
	"pos-restaurant/models"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(ctx context.Context, c *models.Customer) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO customers (hotel_guest_id, type, name, phone, visit_count, last_visit)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING cust_id
	`, c.HotelGuestID, c.Type, c.Name, c.Phone, c.VisitCount, c.LastVisit).Scan(&id)

	return id, err
}

func (r *CustomerRepository) List(ctx context.Context) ([]*models.Customer, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT cust_id, hotel_guest_id, type, name, phone, visit_count, last_visit
		FROM customers
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		var c models.Customer
		err := rows.Scan(&c.CustID, &c.HotelGuestID, &c.Type, &c.Name, &c.Phone, &c.VisitCount, &c.LastVisit)
		if err != nil {
			return nil, err
		}
		customers = append(customers, &c)
	}
	return customers, nil
}

func (r *CustomerRepository) GetByID(ctx context.Context, id int) (*models.Customer, error) {
	var c models.Customer
	err := r.db.QueryRowContext(ctx, `
		SELECT cust_id, hotel_guest_id, type, name, phone, visit_count, last_visit
		FROM customers
		WHERE cust_id = $1
	`, id).Scan(&c.CustID, &c.HotelGuestID, &c.Type, &c.Name, &c.Phone, &c.VisitCount, &c.LastVisit)

	return &c, err
}

func (r *CustomerRepository) Update(ctx context.Context, c *models.Customer) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE customers SET hotel_guest_id = $1, type = $2, name = $3,
		phone = $4, visit_count = $5, last_visit = $6, updated_at = NOW()
		WHERE cust_id = $7
	`, c.HotelGuestID, c.Type, c.Name, c.Phone, c.VisitCount, c.LastVisit, c.CustID)
	return err
}

func (r *CustomerRepository) SoftDelete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE customers SET updated_at = NULL WHERE cust_id = $1
	`, id)
	return err
}
