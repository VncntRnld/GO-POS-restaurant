package models

import (
	"database/sql"
	"time"
)

// Reservations
type Reservation struct {
	ID              int            `json:"id"`
	CustomerID      int            `json:"customer_id"`
	ReservationTime time.Time      `json:"reservation_time"`
	Pax             int            `json:"pax"`
	TableID         int            `json:"table_id"`
	Status          string         `json:"status"`
	SpecialRequest  sql.NullString `json:"special_request"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type ReservationWithDetails struct {
	ID              int            `json:"id"`
	CustomerName    string         `json:"customer_name"`
	ReservationTime time.Time      `json:"reservation_time"`
	Pax             int            `json:"pax"`
	TableNumber     string         `json:"table_number"`
	Status          string         `json:"status"`
	SpecialRequest  sql.NullString `json:"special_request"`
}
