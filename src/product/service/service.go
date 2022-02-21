package service

import (
	"context"

	"github.com/corneliusdavid97/laundry-go/src/product"
)

type Service struct {
	store Store
}

type Store interface {
	GetAllProduct(ctx context.Context, filter product.Filter) ([]product.Product, error)
}

func (s *Service) GetAllActiveProducts(ctx context.Context, filter product.Filter) ([]product.Product, error) {
	res, err := s.store.GetAllProduct(ctx, filter)
	if err != nil {
		return []product.Product{}, err
	}
	return res, nil
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}
