package store

import (
	"context"
	"log"

	"github.com/corneliusdavid97/laundry-go/src/customer"
	"github.com/jmoiron/sqlx"
)

const queryGetAllCustomer = `
	select
		id,
		name,
		coalesce(phone,''),
		coalesce(address,''),
		active
	from
		cust_data
	where
		active=$1
`
const queryGetCustomerByID = `
	select
		id,
		name,
		coalesce(phone,''),
		coalesce(address,''),
		active
	from
		cust_data
	where
		id=$1
`

const queryInsertNewCustomer = `
	insert into cust_data (
		name, 
		phone, 
		address
	)values(
		$1,
		$2,
		$3
	)
	
`

type Store struct {
	getDB func(dbName, replication string) (*sqlx.DB, error)
}

func (s *Store) GetAllCustomer(ctx context.Context, active bool) ([]customer.Customer, error) {
	db, err := s.getDB("db_main", "master")
	if err != nil {
		return []customer.Customer{}, err
	}

	rows, err := db.QueryContext(ctx, queryGetAllCustomer, active)
	if err != nil {
		return []customer.Customer{}, err
	}
	var res []customer.Customer
	for rows.Next() {
		var cust customer.Customer
		err = rows.Scan(&cust.ID, &cust.Name, &cust.PhoneNumber, &cust.Address, &cust.Active)
		if err != nil {
			log.Printf("Failed to scan customer, err:%v, cust:%v", err, cust)
		} else {
			res = append(res, cust)
		}
	}
	return res, nil
}

func (s *Store) GetCustomerByID(ctx context.Context, ID int64) (customer.Customer, error) {
	db, err := s.getDB("db_main", "master")
	if err != nil {
		return customer.Customer{}, err
	}
	row := db.QueryRowContext(ctx, queryGetCustomerByID, ID)
	var cust customer.Customer
	err = row.Scan(&cust.ID, &cust.Name, &cust.PhoneNumber, &cust.Address, &cust.Active)
	if err != nil {
		return customer.Customer{}, err
	}
	return cust, nil
}

func (s *Store) InsertNewCustomer(ctx context.Context, cust customer.Customer) error {
	db, err := s.getDB("db_main", "master")
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, queryInsertNewCustomer, cust.Name, cust.PhoneNumber, cust.Address)
	if err != nil {
		return err
	}
	return nil
}

func NewStore(getDB func(dbName, replication string) (*sqlx.DB, error)) *Store {
	return &Store{
		getDB: getDB,
	}
}
