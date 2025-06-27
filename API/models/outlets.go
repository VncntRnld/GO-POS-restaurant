package models

import (
	"database/sql"
	"time"
)

// Outlets
type Outlet struct {
	ID                   int            `json:"id"`
	Name                 string         `json:"name"`
	Location             sql.NullString `json:"location"`
	ServiceChargePercent float64        `json:"service_charge_percentage"`
	TaxPercentage        float64        `json:"tax_percentage"`
	IsActive             bool           `json:"is_active"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            sql.NullTime   `json:"deleted_at"`
}

// Tables
type Table struct {
	ID           int            `json:"id"`
	OutletID     int            `json:"outlet_id"`
	TableNumber  string         `json:"table_number"`
	Capacity     int            `json:"capacity"`
	LocationType sql.NullString `json:"location_type"` // 'Indoor' or 'Outdoor'
	Status       string         `json:"status"`        // available, occupied, etc.
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    sql.NullTime   `json:"deleted_at"`
}

// Staff
type Staff struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Role      string       `json:"role"`
	PinCode   string       `json:"pin_code"`
	IsActive  bool         `json:"is_active"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}
