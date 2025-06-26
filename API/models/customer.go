package models

import (
	"database/sql"
	"time"
)

// Customers
type Customer struct {
	CustID       int            `json:"cust_id"`
	HotelGuestID sql.NullString `json:"hotel_guest_id"`
	Tipe         string         `json:"tipe"`
	Nama         string         `json:"nama"`
	Phone        sql.NullString `json:"phone"`
	VisitCount   int            `json:"visit_count"`
	LastVisit    sql.NullTime   `json:"last_visit"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// Customer Visits
type CustomerVisit struct {
	ID            int             `json:"id"`
	CustomerID    int             `json:"customer_id"`
	VisitType     string          `json:"visit_type"`
	VisitDate     time.Time       `json:"visit_date"`
	RoomNumber    sql.NullString  `json:"room_number"`
	ReservationID sql.NullInt64   `json:"reservation_id"`
	OutletID      int             `json:"outlet_id"`
	TotalSpent    sql.NullFloat64 `json:"total_spent"`
	Pax           int             `json:"pax"`
	CreatedAt     time.Time       `json:"created_at"`
}
