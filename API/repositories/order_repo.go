package repositories

import (
	"context"
	"database/sql"
	"log"
	"pos-restaurant/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, req *models.OrderRequest) (int, error) {
	log.Println("⏳ Begin transaction...")
	tx, err := r.db.BeginTx(ctx, nil)

	defer func() {
		if p := recover(); p != nil {
			log.Println("❌ Panic recovered, rollback transaction")
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Printf("❌ Error occurred (%v), rollback transaction", err)
			tx.Rollback()
		}
	}()

	// Insert ke tabel orders
	log.Println("➡️ Inserting into orders table...")
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
		log.Printf("❌ Failed to insert order: %v", err)
		return 0, err
	}
	log.Printf("✅ Order inserted with ID: %d", orderID)

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
			log.Printf("❌ Failed to insert order item: %v", err)
			return 0, err
		}
		log.Printf("✅ Order item inserted with ID: %d", orderItemID)

		// Simpan excluded ingredients
		excludedMap := make(map[int]bool)
		for _, ingID := range item.ExcludedIngredientIDs {
			excludedMap[ingID] = true

			log.Printf("➡️ Inserting excluded ingredient %d...", ingID)
			_, err = tx.ExecContext(ctx, `
				INSERT INTO order_item_ingredient_excluded (
					order_item_id, ingredient_id
				) VALUES ($1, $2)
			`, orderItemID, ingID)
			if err != nil {
				log.Printf("❌ Failed to insert excluded ingredient: %v", err)
				return 0, err
			}
			log.Println("✅ Excluded ingredient inserted")
		}

		// // Ambil bahan dari menu item
		// log.Printf("➡️ Fetching ingredients for menu_item_id %d...", item.MenuItemID)
		// rows, err := tx.QueryContext(ctx, `
		// 	SELECT ingredient_id, qty
		// 	FROM menu_ingredients
		// 	WHERE menu_item_id = $1
		// `, item.MenuItemID)
		// if err != nil {
		// 	log.Printf("❌ Failed to query menu ingredients: %v", err)
		// 	return 0, err
		// }
		// defer rows.Close()

		// for rows.Next() {
		// 	var ingredientID int
		// 	var usedQty float64
		// 	if err := rows.Scan(&ingredientID, &usedQty); err != nil {
		// 		log.Printf("❌ Failed to scan menu ingredient row: %v", err)
		// 		return 0, err
		// 	}

		// 	if excludedMap[ingredientID] {
		// 		log.Printf("⚠️ Ingredient %d is excluded, skipping...", ingredientID)
		// 		continue
		// 	}

		// 	totalUsed := usedQty * float64(item.Qty)
		// 	log.Printf("➡️ Validating stock for ingredient %d (need %.2f units)...", ingredientID, totalUsed)

		// 	var currentQty float64
		// 	err = tx.QueryRowContext(ctx, `
		// 		SELECT qty FROM ingredients WHERE id = $1
		// 	`, ingredientID).Scan(&currentQty)
		// 	if err != nil {
		// 		log.Printf("❌ Failed to get ingredient stock: %v", err)
		// 		return 0, err
		// 	}
		// 	if currentQty < totalUsed {
		// 		log.Printf("❌ Not enough stock for ingredient %d (available %.2f, needed %.2f)", ingredientID, currentQty, totalUsed)
		// 		return 0, fmt.Errorf("stok bahan %d tidak cukup", ingredientID)
		// 	}

		// 	log.Printf("➡️ Updating stock for ingredient %d: reducing by %.2f", ingredientID, totalUsed)
		// 	_, err := tx.ExecContext(ctx, `
		// 		UPDATE ingredients
		// 		SET qty = qty - $1
		// 		WHERE id = $2
		// 	`, totalUsed, ingredientID)
		// 	if err != nil {
		// 		log.Printf("❌ Failed to update stock: %v", err)
		// 		return 0, err
		// 	}
		// 	log.Printf("✅ Stock updated for ingredient %d", ingredientID)
		// }
	}

	log.Println("💾 Committing transaction...")
	if err := tx.Commit(); err != nil {
		log.Printf("❌ Failed to commit transaction: %v", err)
		return 0, err
	}
	log.Printf("✅ Order transaction committed successfully with ID %d", orderID)

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
				Qty:                   float64(qty),
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
				Qty:                   float64(qty),
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

func (r *OrderRepository) SoftDelete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE orders SET status = 'void', updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}
