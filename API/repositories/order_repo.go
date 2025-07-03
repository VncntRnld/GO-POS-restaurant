package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"pos-restaurant/models"
	"slices"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

type IngredientUsage struct {
	IngredientID int
	UsedQty      float64
}

func (r *OrderRepository) Create(ctx context.Context, req *models.OrderRequest) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// Insert ke tabel orders
	var orderID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO orders (
			order_number, table_id, customer_id, hotel_room,
			waiter_id, outlet_id, status, order_type
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id
	`, req.OrderNumber, req.TableID, req.CustomerID, req.HotelRoom,
		req.WaiterID, req.OutletID, req.Status, req.OrderType,
	).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	// Masukkan menu berdasarkan order
	for _, item := range req.Items {
		log.Println("➡️ Inserting order item...")
		var orderItemID int
		err = tx.QueryRowContext(ctx, `
			INSERT INTO order_items (
				order_id, menu_item_id, qty, notes, unit_price
			) VALUES ($1,$2,$3,$4,$5)
			RETURNING id
		`, orderID, item.MenuItemID, item.Qty, item.Notes, item.UnitPrice).Scan(&orderItemID)
		if err != nil {
			return 0, err
		}

		// Simpan excluded ingredients
		excludedMap := make(map[int]bool)
		for _, ingID := range item.ExcludedIngredientIDs {
			excludedMap[ingID] = true

			_, err = tx.ExecContext(ctx, `
				INSERT INTO order_item_ingredient_excluded (
					order_item_id, ingredient_id
				) VALUES ($1, $2)
			`, orderItemID, ingID)
			if err != nil {
				return 0, err
			}
		}

		// Ambil bahan dari menu item
		rows, err := tx.QueryContext(ctx, `
			SELECT ingredient_id, qty
			FROM menu_ingredients
			WHERE menu_item_id = $1
		`, item.MenuItemID)
		if err != nil {
			return 0, err
		}

		var ingredients []IngredientUsage
		for rows.Next() {
			var ing IngredientUsage
			if err := rows.Scan(&ing.IngredientID, &ing.UsedQty); err != nil {
				rows.Close()
				return 0, err
			}
			ingredients = append(ingredients, ing)
		}
		rows.Close()

		// Proses setiap bahan
		for _, ing := range ingredients {
			if excludedMap[ing.IngredientID] {
				continue
			}

			totalUsed := ing.UsedQty * float64(item.Qty)

			var currentQty float64
			err = tx.QueryRowContext(ctx, `
			SELECT qty FROM ingredients WHERE id = $1
		`, ing.IngredientID).Scan(&currentQty)
			if err != nil {
				return 0, err
			}

			if currentQty < totalUsed {
				return 0, fmt.Errorf("stok bahan %d tidak cukup", ing.IngredientID)
			}

			_, err = tx.ExecContext(ctx, `
			UPDATE ingredients
			SET qty = qty - $1
			WHERE id = $2
		`, totalUsed, ing.IngredientID)
			if err != nil {
				return 0, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return orderID, nil
}

func (r *OrderRepository) List(ctx context.Context) ([]*models.OrderRequest, error) {
	query := `
	SELECT 
		o.id, o.order_number, o.table_id, o.customer_id, o.hotel_room,
		o.waiter_id, o.outlet_id, o.status, o.order_type,
		oi.id, oi.menu_item_id, oi.qty, oi.notes, oi.unit_price,
		ie.ingredient_id
	FROM orders o
	LEFT JOIN order_items oi ON o.id = oi.order_id
	LEFT JOIN order_item_ingredient_excluded ie ON oi.id = ie.order_item_id
	ORDER BY o.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderMap := map[int]*models.OrderRequest{}

	for rows.Next() {
		var (
			orderID, tableID, customerID, waiterID, outletID int
			orderNumber, status, orderType                   string
			orderItemID, menuItemID                          int
			UnitPrice, qty                                   float64
			hotelRoom, notes                                 sql.NullString
			excludedIngID                                    sql.NullInt64
		)

		err := rows.Scan(
			&orderID, &orderNumber, &tableID, &customerID, &hotelRoom,
			&waiterID, &outletID, &status, &orderType,
			&orderItemID, &menuItemID, &qty, &notes, &UnitPrice,
			&excludedIngID,
		)
		if err != nil {
			return nil, err
		}

		order, exists := orderMap[orderID]
		if !exists {
			order = &models.OrderRequest{
				ID:          orderID,
				OrderNumber: orderNumber,
				TableID:     tableID,
				CustomerID:  customerID,
				HotelRoom:   hotelRoom,
				WaiterID:    waiterID,
				OutletID:    outletID,
				Status:      status,
				OrderType:   orderType,
				Items:       []models.OrderItemInput{},
			}
			orderMap[orderID] = order
		}

		if orderItemID != 0 {
			item := &models.OrderItemInput{
				ID:                    orderItemID,
				MenuItemID:            menuItemID,
				Qty:                   qty,
				Notes:                 notes.String,
				UnitPrice:             UnitPrice,
				ExcludedIngredientIDs: []int{},
			}
			if excludedIngID.Valid {
				item.ExcludedIngredientIDs = append(item.ExcludedIngredientIDs, int(excludedIngID.Int64))
			}
			order.Items = append(order.Items, *item)
		}
	}

	var results []*models.OrderRequest
	for _, order := range orderMap {
		results = append(results, order)
	}

	return results, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id int) (*models.OrderRequest, error) {
	query := `
	SELECT 
		o.id, o.order_number, o.table_id, o.customer_id, o.hotel_room,
		o.waiter_id, o.outlet_id, o.status, o.order_type,
		oi.id, oi.menu_item_id, oi.qty, oi.notes, oi.unit_price,
		ie.ingredient_id
	FROM orders o
	LEFT JOIN order_items oi ON o.id = oi.order_id
	LEFT JOIN order_item_ingredient_excluded ie ON oi.id = ie.order_item_id
	WHERE o.id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var order *models.OrderRequest

	for rows.Next() {
		var (
			orderID, tableID, customerID, waiterID, outletID int
			orderNumber, status, orderType                   string
			orderItemID, menuItemID                          int
			UnitPrice, qty                                   float64
			hotelRoom, notes                                 sql.NullString
			excludedIngID                                    sql.NullInt64
		)

		err := rows.Scan(
			&orderID, &orderNumber, &tableID, &customerID, &hotelRoom,
			&waiterID, &outletID, &status, &orderType,
			&orderItemID, &menuItemID, &qty, &notes, &UnitPrice,
			&excludedIngID,
		)
		if err != nil {
			return nil, err
		}

		if order == nil {
			order = &models.OrderRequest{
				ID:          orderID,
				OrderNumber: orderNumber,
				TableID:     tableID,
				CustomerID:  customerID,
				HotelRoom:   hotelRoom,
				WaiterID:    waiterID,
				OutletID:    outletID,
				Status:      status,
				OrderType:   orderType,
				Items:       []models.OrderItemInput{},
			}
		}

		if orderItemID != 0 {
			item := &models.OrderItemInput{
				ID:                    orderItemID,
				MenuItemID:            menuItemID,
				Qty:                   qty,
				Notes:                 notes.String,
				UnitPrice:             UnitPrice,
				ExcludedIngredientIDs: []int{},
			}
			if excludedIngID.Valid {
				item.ExcludedIngredientIDs = append(item.ExcludedIngredientIDs, int(excludedIngID.Int64))
			}
			order.Items = append(order.Items, *item)
		}
	}

	if order == nil {
		return nil, sql.ErrNoRows
	}

	return order, nil
}

func (r *OrderRepository) Update(ctx context.Context, order *models.Order) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE orders SET
			table_id = $1,
			customer_id = $2,
			hotel_room = $3,
			waiter_id = $4,
			outlet_id = $5,
			status = $6,
			order_type = $7,
			updated_at = NOW()
		WHERE id = $8
	`,
		order.TableID,
		order.CustomerID,
		order.HotelRoom,
		order.WaiterID,
		order.OutletID,
		order.Status,
		order.OrderType,
		order.ID,
	)
	return err
}

func (r *OrderRepository) AddItem(ctx context.Context, orderID int, item *models.AddOrderItemRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// 1. Tambahkan order item
	var orderItemID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO order_items (
			order_id, menu_item_id, qty, notes, unit_price
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, orderID, item.MenuItemID, item.Qty, item.Notes, item.UnitPrice).Scan(&orderItemID)
	if err != nil {
		return err
	}

	// 2. Tambahkan excluded ingredients (jika ada)
	for _, ingID := range item.ExcludedIngredientIDs {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO order_item_ingredient_excluded (
				order_item_id, ingredient_id
			) VALUES ($1, $2)
		`, orderItemID, ingID)
		if err != nil {
			return err
		}
	}

	// 3. Ambil bahan dari menu_ingredients dan simpan ke slice
	var ingredients []IngredientUsage

	rows, err := tx.QueryContext(ctx, `
		SELECT ingredient_id, qty
		FROM menu_ingredients
		WHERE menu_item_id = $1
	`, item.MenuItemID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var ing IngredientUsage
		if err := rows.Scan(&ing.IngredientID, &ing.UsedQty); err != nil {
			rows.Close()
			return err
		}
		ingredients = append(ingredients, ing)
	}
	rows.Close()

	// 4. Proses stok dan update ingredients
	for _, ing := range ingredients {
		// Skip jika termasuk dalam excluded
		if slices.Contains(item.ExcludedIngredientIDs, ing.IngredientID) {
			continue
		}

		// Ambil stok
		var currentQty float64
		err := tx.QueryRowContext(ctx, `
			SELECT qty FROM ingredients WHERE id = $1
		`, ing.IngredientID).Scan(&currentQty)
		if err != nil {
			return err
		}

		totalNeeded := ing.UsedQty * item.Qty
		if currentQty < totalNeeded {
			return fmt.Errorf("stok tidak cukup untuk bahan id %d", ing.IngredientID)
		}

		// Kurangi stok
		_, err = tx.ExecContext(ctx, `
			UPDATE ingredients
			SET qty = qty - $1
			WHERE id = $2
		`, totalNeeded, ing.IngredientID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) SoftDelete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE orders SET status = 'void', updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}
