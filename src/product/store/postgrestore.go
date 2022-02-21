package store

import (
	"context"
	"fmt"
	"log"

	"github.com/corneliusdavid97/laundry-go/src/product"
	"github.com/jmoiron/sqlx"
)

const queryGetAllProduct = `
	select
		id,
		product_name,
		price_standard,
		price_express_today,
		price_express_tmr,
		active,
		is_satuan
	from
		product_data
	where
		%s
`

type Store struct {
	getDB func(dbName, replication string) (*sqlx.DB, error)
}

func (s *Store) GetAllProduct(ctx context.Context, filter product.Filter) ([]product.Product, error) {
	db, err := s.getDB("db_main", "master")
	if err != nil {
		return []product.Product{}, err
	}
	query := fmt.Sprintf(queryGetAllProduct, constructFilter(filter))
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return []product.Product{}, err
	}
	res := make([]product.Product, 0)
	for rows.Next() {
		var p product.Product
		err = rows.Scan(&p.ID, &p.Name, &p.PriceStandard, &p.PriceExpressToday, &p.PriceExpressTomorrow, &p.Active, &p.IsSatuan)
		if err != nil {
			log.Printf("Failed to scan product, err:%v, product:%v", err, p)
		}
		res = append(res, p)
	}
	return res, nil
}

func constructFilter(filter product.Filter) string {
	if (filter == product.Filter{}) {
		return "active=true"
	}
	var filters []string
	if filter.Active != nil {
		filters = append(filters, fmt.Sprintf("active=%v", *filter.Active))
	}
	if filter.IsSatuan != nil {
		filters = append(filters, fmt.Sprintf("is_satuan=%v", *filter.IsSatuan))
	}

	res := ""
	for i, f := range filters {
		if i != 0 {
			res += "and "
		}
		res += f + " "
	}
	return res
}

func NewStore(getDB func(dbName, replication string) (*sqlx.DB, error)) *Store {
	return &Store{
		getDB: getDB,
	}
}
