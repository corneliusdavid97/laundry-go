package service

import (
	"context"

	"github.com/corneliusdavid97/laundry-go/src/transaction"
)

type Service struct {
	store Store
}

type Store interface {
	NewTransaction(ctx context.Context, trans transaction.Transaction) error
	MarkDateTaken(ctx context.Context, ID int64) error
	GetTransactionDataByID(ctx context.Context, ID int64) (transaction.Transaction, error)
}

func (s *Service) GetTransactionDataByID(ctx context.Context, ID int64) (transaction.Transaction, error) {
	res, err := s.store.GetTransactionDataByID(ctx, ID)
	if err != nil {
		return transaction.Transaction{}, err
	}
	return res, nil
}

func (s *Service) MarkDateTaken(ctx context.Context, ID int64) error {
	err := s.store.MarkDateTaken(ctx, ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) NewTransaction(ctx context.Context, trans transaction.Transaction) error {
	err := s.store.NewTransaction(ctx, trans)
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
