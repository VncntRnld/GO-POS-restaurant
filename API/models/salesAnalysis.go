package models

import "time"

// Sales Analysis Daily
type SalesAnalysisDaily struct {
	ID               int       `json:"id"`
	OutletID         int       `json:"outlet_id"`
	AnalysisDate     time.Time `json:"analysis_date"`
	TotalSales       float64   `json:"total_sales"`
	TotalCovers      int       `json:"total_covers"`
	AvgSpendPerCover float64   `json:"avg_spend_per_cover"`
	DiscountAmount   float64   `json:"discount_amount"`
	VoidAmount       float64   `json:"void_amount"`
	CreatedAt        time.Time `json:"created_at"`
}
