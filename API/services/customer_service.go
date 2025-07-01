// services/customer_service.go
package services

import (
	"context"
	"pos-restaurant/models"
	"pos-restaurant/repositories"
)

type CustomerService struct {
	repo *repositories.CustomerRepository
}

func NewCustomerService(repo *repositories.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, c *models.Customer) (int, error) {
	return s.repo.Create(ctx, c)
}

func (s *CustomerService) GetAllCustomers(ctx context.Context) ([]*models.Customer, error) {
	return s.repo.List(ctx)
}

func (s *CustomerService) GetCustomerByID(ctx context.Context, id int) (*models.Customer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, c *models.Customer) error {
	return s.repo.Update(ctx, c)
}

func (s *CustomerService) SoftDeleteCustomer(ctx context.Context, id int) error {
	return s.repo.SoftDelete(ctx, id)
}
