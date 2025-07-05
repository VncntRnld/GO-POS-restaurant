package repositories

import (
	"context"
	"database/sql"
	"pos-restaurant/models"
)

type TableTransferRepository struct {
	db *sql.DB
}

func NewTableTransferRepository(db *sql.DB) *TableTransferRepository {
	return &TableTransferRepository{db}
}

func (r *TableTransferRepository) Create(ctx context.Context, t *models.TableTransfer) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO table_transfers (
			order_id, from_table_id, to_table_id, transferred_by, reason
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, t.OrderID, t.FromTableID, t.ToTableID, t.TransferredBy, t.Reason).Scan(&id)
	if err != nil {
		return 0, err
	}

	// 🔁 Update table_id di orders
	_, err = tx.ExecContext(ctx, `
		UPDATE orders
		SET table_id = $1,
		updated_at = NOW()
		WHERE id = $2
	`, t.ToTableID, t.OrderID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TableTransferRepository) List(ctx context.Context) ([]models.TableTransfer, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, order_id, from_table_id, to_table_id, transferred_by, transferred_at, reason
		FROM table_transfers
		ORDER BY transferred_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.TableTransfer
	for rows.Next() {
		var t models.TableTransfer
		err := rows.Scan(&t.ID, &t.OrderID, &t.FromTableID, &t.ToTableID, &t.TransferredBy, &t.TransferredAt, &t.Reason)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *TableTransferRepository) GetByID(ctx context.Context, id int) (*models.TableTransfer, error) {
	var t models.TableTransfer
	err := r.db.QueryRowContext(ctx, `
		SELECT id, order_id, from_table_id, to_table_id, transferred_by, transferred_at, reason
		FROM table_transfers
		WHERE id = $1
	`, id).Scan(&t.ID, &t.OrderID, &t.FromTableID, &t.ToTableID, &t.TransferredBy, &t.TransferredAt, &t.Reason)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TableTransferRepository) Update(ctx context.Context, t *models.TableTransfer) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE table_transfers
		SET order_id = $1, from_table_id = $2, to_table_id = $3,
		    transferred_by = $4, reason = $5
		WHERE id = $6
	`, t.OrderID, t.FromTableID, t.ToTableID, t.TransferredBy, t.Reason, t.ID)
	if err != nil {
		return err
	}

	// 🔁 Update table_id di orders
	_, err = tx.ExecContext(ctx, `
		UPDATE orders
		SET table_id = $1,
		updated_at = NOW()
		WHERE id = $2
	`, t.ToTableID, t.OrderID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *TableTransferRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM table_transfers WHERE id = $1`, id)
	return err
}
