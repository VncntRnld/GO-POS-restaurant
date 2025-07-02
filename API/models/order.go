package models

import (
	"database/sql"
	"time"
)

// Orders
type Order struct {
	ID          int            `json:"id"`
	OrderNumber string         `json:"order_number"`
	TableID     int            `json:"table_id"`
	CustomerID  int            `json:"customer_id"`
	HotelRoom   sql.NullString `json:"hotel_room"`
	WaiterID    int            `json:"waiter_id"`
	OutletID    int            `json:"outlet_id"`
	Status      string         `json:"status"`
	OrderType   string         `json:"order_type"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type OrderRequest struct {
	ID          int              `json:"id"`
	OrderNumber string           `json:"order_number"`
	TableID     int              `json:"table_id"`
	CustomerID  int              `json:"customer_id"`
	HotelRoom   sql.NullString   `json:"hotel_room"`
	WaiterID    int              `json:"waiter_id"`
	OutletID    int              `json:"outlet_id"`
	Status      string           `json:"status"`
	OrderType   string           `json:"order_type"`
	Items       []OrderItemInput `json:"items"`
}

type OrderItemInput struct {
	ID                    int     `json:"id"`
	MenuItemID            int     `json:"menu_item_id"`
	Qty                   float64 `json:"qty"`
	Notes                 string  `json:"notes,omitempty"`
	UnitPrice             float64 `json:"unit_price"` // captured per item
	ExcludedIngredientIDs []int   `json:"excluded_ingredients"`
}

// Bills
type Bill struct {
	ID             int           `json:"id"`
	BillNumber     string        `json:"bill_number"`
	OrderID        int           `json:"order_id"`
	OriginalBillID sql.NullInt64 `json:"original_bill_id"`
	Status         string        `json:"status"`
	Subtotal       float64       `json:"subtotal"`
	TaxAmount      float64       `json:"tax_amount"`
	ServiceCharge  float64       `json:"service_charge"`
	DiscountAmount float64       `json:"discount_amount"`
	TotalAmount    float64       `json:"total_amount"`
	PaidAmount     float64       `json:"paid_amount"`
	BalanceDue     float64       `json:"balance_due"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// Bill Payments
type BillPayment struct {
	ID                   int            `json:"id"`
	BillID               int            `json:"bill_id"`
	PaymentMethod        string         `json:"payment_method"`
	Amount               float64        `json:"amount"`
	ReferenceNumber      sql.NullString `json:"reference_number"`
	RoomChargeApprovedBy sql.NullInt64  `json:"room_charge_approved_by"`
	PaymentTime          time.Time      `json:"payment_time"`
}

// Table Transfers
type TableTransfer struct {
	ID            int            `json:"id"`
	OrderID       int            `json:"order_id"`
	FromTableID   int            `json:"from_table_id"`
	ToTableID     int            `json:"to_table_id"`
	TransferredBy int            `json:"transferred_by"`
	TransferredAt time.Time      `json:"transferred_at"`
	Reason        sql.NullString `json:"reason"`
}
