package product

import "context"

type Product struct {
	ID                   int64   `json:"id"`
	Name                 string  `json:"name"`
	PriceStandard        float64 `json:"price_standard"`
	PriceExpressToday    float64 `json:"price_express_today"`
	PriceExpressTomorrow float64 `json:"price_express_tomorrow"`
	IsSatuan             bool    `json:"is_satuan"`
	Active               bool    `json:"active"`
}

type Filter struct {
	IsSatuan *bool
	Active   *bool
}

type Service interface {
	GetAllActiveProducts(ctx context.Context, filter Filter) ([]Product, error)
	// GetProductByID(ctx context.Context, ID int64) (Product, error)
	// AddNewProduct(ctx context.Context, product Product) error
}

var defaultService Service

func Init(s Service) {
	defaultService = s
}

func GetService() Service {
	return defaultService
}
