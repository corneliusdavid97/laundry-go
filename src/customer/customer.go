package customer

import (
	"context"
	"errors"
)

type Customer struct {
	ID          int64
	Name        string
	PhoneNumber string
	Address     string
	Active      bool
}

var ErrInvalidCustomer = errors.New("Invalid customer data")

type Service interface {
	GetAllActiveCustomer(ctx context.Context) ([]Customer, error)
	GetCustomerByID(ctx context.Context, id int64) (Customer, error)
	InsertNewCustomer(ctx context.Context, cust Customer) error
}

var defaultService Service

func Init(s Service) {
	defaultService = s
}

func GetService() Service {
	return defaultService
}
