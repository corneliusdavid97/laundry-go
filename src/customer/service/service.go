package service

import (
	"context"

	"github.com/corneliusdavid97/laundry-go/src/customer"
)

type Service struct {
	store Store
}

type Store interface {
	GetAllCustomer(ctx context.Context, active bool) ([]customer.Customer, error)
	GetCustomerByID(ctx context.Context, ID int64) (customer.Customer, error)
	InsertNewCustomer(ctx context.Context, cust customer.Customer) error
}

func (s *Service) GetAllActiveCustomer(ctx context.Context) ([]customer.Customer, error) {
	custs, err := s.store.GetAllCustomer(ctx, true)
	if err != nil {
		return []customer.Customer{}, err
	}
	return custs, nil
}

func (s *Service) GetCustomerByID(ctx context.Context, ID int64) (customer.Customer, error) {
	cust, err := s.store.GetCustomerByID(ctx, ID)
	if err != nil {
		return customer.Customer{}, err
	}
	return cust, nil
}

func (s *Service) InsertNewCustomer(ctx context.Context, cust customer.Customer) error {
	if len(cust.Name) == 0 {
		return customer.ErrInvalidCustomer
	}
	err := s.store.InsertNewCustomer(ctx, cust)
	if err != nil {
		return err
	}
	return nil
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}
