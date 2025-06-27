package models

import (
	"database/sql"
	"time"
)

// Menu Category
type MenuCategory struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	CreatedAt time.Time    `json:"created_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

// Ingredients
type Ingredient struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Qty         float64        `json:"qty"`
	Unit        string         `json:"unit"`
	IsAllergen  bool           `json:"is_allergen"`
	IsActive    bool           `json:"is_active"`
	Description sql.NullString `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   sql.NullTime   `json:"deleted_at"`
}

// Menu Items
type MenuItem struct {
	ID              int            `json:"id"`
	CategoryID      int            `json:"category_id"`
	SKU             string         `json:"sku"`
	Name            string         `json:"name"`
	Description     sql.NullString `json:"description"`
	Price           float64        `json:"price"`
	Cost            float64        `json:"cost"`
	IsActive        bool           `json:"is_active"`
	PreparationTime sql.NullInt64  `json:"preparation_time"`
	Tags            []string       `json:"tags"` // parsed manually from JSONB
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       sql.NullTime   `json:"deleted_at"`
}

// Menu Ingredients
type MenuIngredient struct {
	ID           int       `json:"id"`
	MenuItemID   int       `json:"menu_item_id"`
	IngredientID int       `json:"ingredient_id"`
	Qty          float64   `json:"qty"`
	IsRemovable  bool      `json:"is_removable"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Menu details with ingredients
type MenuItemWithIngredients struct {
	MenuItem
	Ingredients []Ingredient `json:"ingredients"`
}
