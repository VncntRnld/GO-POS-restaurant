package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"pos-restaurant/models"

	"github.com/google/uuid"
)

type BillRepository struct {
	db *sql.DB
}

func NewBillRepository(db *sql.DB) *BillRepository {
	return &BillRepository{db: db}
}

func (r *BillRepository) Create(ctx context.Context, orderID int, discount float64) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// 1. Ambil item dan outlet dari order
	var outletID int
	rows, err := tx.QueryContext(ctx, `
		SELECT oi.qty, oi.unit_price, o.outlet_id
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		WHERE o.id = $1
	`, orderID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var subtotal float64
	for rows.Next() {
		var qty, unitPrice float64
		if err := rows.Scan(&qty, &unitPrice, &outletID); err != nil {
			return 0, err
		}
		subtotal += qty * unitPrice
	}

	// 2. Ambil tax & service dari outlet
	var taxPct, servicePct float64
	err = tx.QueryRowContext(ctx, `
		SELECT tax_percentage, service_charge_percentage
		FROM outlets
		WHERE id = $1
	`, outletID).Scan(&taxPct, &servicePct)
	if err != nil {
		return 0, err
	}

	serviceCharge := subtotal * servicePct / 100
	taxAmount := (subtotal + serviceCharge) * taxPct / 100
	total := subtotal + serviceCharge + taxAmount - discount
	if total < 0 {
		total = 0
	}

	// 3. Insert ke Bill
	var billID int
	billNumber := uuid.New().String()
	err = tx.QueryRowContext(ctx, `
		INSERT INTO bills (
			bill_number, order_id, status, subtotal,
			tax_amount, service_charge, discount_amount,
			total_amount
		) VALUES (
			$1, $2, 'open', $3, $4, $5, $6, $7)
		RETURNING id
	`, billNumber, orderID, subtotal, taxAmount, serviceCharge, discount, total).Scan(&billID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return billID, nil
}

func (r *BillRepository) CreateSplit(ctx context.Context, req models.SplitBillRequest) ([]int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 🔒 Validasi item tidak duplikat antar split
	usedItems := make(map[int]bool)
	for _, split := range req.Splits {
		for _, itemID := range split.ItemIDs {
			if usedItems[itemID] {
				return nil, fmt.Errorf("item ID %d digunakan lebih dari satu kali", itemID)
			}
			usedItems[itemID] = true
		}
	}

	var billIDs []int

	for _, split := range req.Splits {
		var subtotal float64
		for _, itemID := range split.ItemIDs {
			var price float64
			err := tx.QueryRowContext(ctx, `SELECT unit_price FROM order_items WHERE id = $1`, itemID).Scan(&price)
			if err != nil {
				return nil, fmt.Errorf("gagal ambil harga untuk item %d: %w", itemID, err)
			}
			subtotal += price
		}

		var outletID int
		err = tx.QueryRowContext(ctx, `
			SELECT o.outlet_id FROM orders o WHERE o.id = $1
		`, req.OriginalOrderID).Scan(&outletID)
		if err != nil {
			return nil, fmt.Errorf("gagal ambil outlet_id dari order %d: %w", req.OriginalOrderID, err)
		}

		var taxPercent, svcPercent float64
		err = tx.QueryRowContext(ctx, `SELECT tax_percentage, service_charge_percentage FROM outlets WHERE id = $1`, outletID).
			Scan(&taxPercent, &svcPercent)
		if err != nil {
			return nil, fmt.Errorf("gagal ambil tax/outlet: %w", err)
		}

		taxAmount := subtotal * taxPercent / 100
		svcAmount := subtotal * svcPercent / 100
		total := subtotal + taxAmount + svcAmount - split.DiscountAmount

		// var paidStatus string
		// switch {
		// case balance <= 0:
		// 	paidStatus = "paid"
		// case split.PaidAmount == 0:
		// 	paidStatus = "unpaid"
		// default:
		// 	paidStatus = "partial"
		// }

		var billID int
		billNumber := uuid.New().String()
		err = tx.QueryRowContext(ctx, `
			INSERT INTO bills (
				bill_number, order_id, original_bill_id, status,
				subtotal, tax_amount, service_charge, discount_amount,
				total_amount
			) VALUES (
				$1, $2, $3, 'open', $4,
				$5, $6, $7, $8
			) RETURNING id
		`,
			billNumber,
			req.OriginalOrderID,
			sql.NullInt64{Int64: int64(req.OriginalBillID), Valid: req.OriginalBillID > 0},
			subtotal, taxAmount, svcAmount,
			split.DiscountAmount, total,
		).Scan(&billID)
		if err != nil {
			return nil, fmt.Errorf("gagal membuat bill: %w", err)
		}

		billIDs = append(billIDs, billID)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return billIDs, nil
}

func (r *BillRepository) GetByID(ctx context.Context, id int) (*models.Bill, error) {
	query := `
		SELECT 
			id, bill_number, order_id, original_bill_id, status,
			subtotal, tax_amount, service_charge, discount_amount,
			total_amount, paid_amount, balance_due,
			created_at, updated_at
		FROM bills
		WHERE id = $1
	`

	var bill models.Bill
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&bill.ID,
		&bill.BillNumber,
		&bill.OrderID,
		&bill.OriginalBillID,
		&bill.Status,
		&bill.Subtotal,
		&bill.TaxAmount,
		&bill.ServiceCharge,
		&bill.DiscountAmount,
		&bill.TotalAmount,
		&bill.PaidAmount,
		&bill.BalanceDue,
		&bill.CreatedAt,
		&bill.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Bisa dikembalikan sebagai not found
		}
		return nil, err
	}

	return &bill, nil
}

func (r *BillRepository) Pay(ctx context.Context, payment *models.BillPayment) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert payment
	_, err = tx.ExecContext(ctx, `
		INSERT INTO bill_payments (
			bill_id, payment_method, amount, reference_number,
			room_charge_approved_by
		) VALUES ($1, $2, $3, $4, $5)
	`,
		payment.BillID,
		payment.PaymentMethod,
		payment.Amount,
		payment.ReferenceNumber,
		payment.RoomChargeApprovedBy,
	)
	if err != nil {
		return err
	}

	// Update bill paid_amount
	_, err = tx.ExecContext(ctx, `
		UPDATE bills
		SET paid_amount = paid_amount + $1
		WHERE id = $2
	`, payment.Amount, payment.BillID)
	if err != nil {
		return err
	}

	// Cek apakah bill sudah lunas
	_, err = tx.ExecContext(ctx, `
		UPDATE bills
		SET 
			status = CASE
				WHEN paid_amount >= total_amount THEN 'paid'
				WHEN paid_amount > 0 AND paid_amount < total_amount THEN 'partial'
				ELSE status
			END
		WHERE id = $1
	`, payment.BillID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
