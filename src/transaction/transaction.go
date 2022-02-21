package transaction

import (
	"context"
	"time"
)

type Transaction struct {
	ID                 int64               `json:"id"`
	CustomerID         int64               `json:"cust_id"`
	GrandTotal         float64             `json:"grand_total"`
	Paid               float64             `json:"paid"`
	TransactionTime    *time.Time          `json:"-"`
	TransactionTimeStr *string             `json:"transaction_time"`
	DueDate            *time.Time          `json:"-"`
	DueDateStr         *string             `json:"due_date"`
	DateTaken          *time.Time          `json:"-"`
	DateTakenStr       *string             `json:"date_taken"`
	PaymentMethod      PaymentMethod       `json:"payment_method"`
	CashierName        string              `json:"cashier_name"`
	Details            []TransactionDetail `json:"details"`
}

type TransactionDetail struct {
	ID          int64   `json:"id"`
	ProductName string  `json:"product_name"`
	ProductType string  `json:"product_type"`
	Price       float64 `json:"price"`
	Quantity    float64 `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

type PaymentMethod string

const (
	PaymentMethodCash      = "cash"
	PaymentMethodShopeePay = "shopee_pay"
	PaymentMethodQRIS      = "qris"
	PaymentMethodBCAMobile = "bca_mobile"
)

type Service interface {
	MarkDateTaken(ctx context.Context, ID int64) error
	NewTransaction(ctx context.Context, trans Transaction) error
	GetTransactionDataByID(ctx context.Context, ID int64) (Transaction, error)
}

var defaultService Service

func Init(s Service) {
	defaultService = s
}

func GetService() Service {
	return defaultService
}
