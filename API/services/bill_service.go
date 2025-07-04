package services

import (
	"context"
	"fmt"
	"pos-restaurant/models"
	"pos-restaurant/repositories"
)

type BillService struct {
	repo *repositories.BillRepository
}

func NewBillService(repo *repositories.BillRepository) *BillService {
	return &BillService{repo: repo}
}

func (s *BillService) Create(ctx context.Context, orderID int, discount float64) (int, error) {
	return s.repo.Create(ctx, orderID, discount)
}

func (s *BillService) CreateSplit(ctx context.Context, req models.SplitBillRequest) ([]int, error) {
	originalBill, err := s.repo.GetByID(ctx, req.OriginalBillID)
	if err != nil {
		return nil, fmt.Errorf("bill tidak ditemukan: %w", err)
	}

	// Validasi bahwa original_order_id di request sama dengan order_id di original bill
	if originalBill.OrderID != req.OriginalOrderID {
		return nil, fmt.Errorf("order_id tidak sesuai dengan original_bill_id yang dituju")
	}

	return s.repo.CreateSplit(ctx, req)
}

func (s *BillService) GetByID(ctx context.Context, id int) (*models.Bill, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BillService) Pay(ctx context.Context, payment *models.BillPayment) error {
	return s.repo.Pay(ctx, payment)
}
