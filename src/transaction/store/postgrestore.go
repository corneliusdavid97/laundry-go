package store

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/corneliusdavid97/laundry-go/src/transaction"
)

const queryInsertTransactionData = `
	insert into transaction_main(
		id,
		customer_id, 
		grand_total, 
		paid,
		due_date,
		payment_method, 
		cashier_name
	)values(
		?,
		?, 
		?, 
		?, 
		?, 
		?,
		?
	)
`

const queryMarkDateTaken = `
	update transaction_main set
		date_taken=now()
	where id=$1
`

const queryInsertTransactionDetail = `
	insert into transaction_detail(
		transaction_id, 
		product_name, 
		product_type, 
		price, 
		quantity, 
		subtotal
	)values
		%s
`

const queryGetTransactionDataByID = `
	select
		id,
		customer_id,
		grand_total,
		paid,
		transaction_time,
		due_date,
		date_taken,
		payment_method,
		cashier_name
	from
		transaction_main
	where
		id=$1	
`

const queryGetTransactionDetailsByTransactionID = `
	select 
		id,
		product_name,
		product_type,
		price,
		quantity,
		subtotal
	from 
		transaction_detail
	where
		transaction_id = $1
`

type Store struct {
	getDB func(dbName, replication string) (*sqlx.DB, error)
}

func (s *Store) GetTransactionDataByID(ctx context.Context, ID int64) (transaction.Transaction, error) {
	db, err := s.getDB("db_main", "master")
	if err != nil {
		return transaction.Transaction{}, err
	}

	tx, err := db.Beginx()
	if err != nil {
		return transaction.Transaction{}, err
	}
	row := tx.QueryRowContext(ctx, queryGetTransactionDataByID, ID)
	var trans transaction.Transaction
	err = row.Scan(&trans.ID, &trans.CustomerID, &trans.GrandTotal, &trans.Paid, &trans.TransactionTime, &trans.DueDate, &trans.DateTaken, &trans.PaymentMethod, &trans.CashierName)
	if err != nil {
		tx.Rollback()
		return transaction.Transaction{}, err
	}
	rows, err := tx.QueryContext(ctx, queryGetTransactionDetailsByTransactionID, ID)
	if err != nil {
		tx.Rollback()
		return transaction.Transaction{}, err
	}

	for rows.Next() {
		var detail transaction.TransactionDetail
		err = rows.Scan(&detail.ID, &detail.ProductName, &detail.ProductType, &detail.Price, &detail.Quantity, &detail.Subtotal)
		if err != nil {
			log.Printf("[Transaction][Store] failed to scan row, detail: %v, err:%v\n", detail, err)
		}
		trans.Details = append(trans.Details, detail)
	}
	tx.Commit()
	return trans, nil
}

func (s *Store) MarkDateTaken(ctx context.Context, ID int64) error {
	db, err := s.getDB("db_main", "master")
	if err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, queryMarkDateTaken, ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) NewTransaction(ctx context.Context, trans transaction.Transaction) error {
	db, err := s.getDB("db_main", "master")
	if err != nil {
		return err
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// insert main data
	query := tx.Rebind(queryInsertTransactionData)
	_, err = tx.ExecContext(ctx, query, trans.ID, trans.CustomerID, trans.GrandTotal, trans.Paid, trans.DueDate, trans.PaymentMethod, trans.CashierName)
	if err != nil {
		tx.Rollback()
		return err
	}

	// insert details
	var params []interface{}
	for _, v := range trans.Details {
		params = append(params, trans.ID, v.ProductName, v.ProductType, v.Price, v.Quantity, v.Subtotal)
	}

	query = tx.Rebind(constructQueryInsertTransactionDetail(len(trans.Details)))
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func constructQueryInsertTransactionDetail(length int) string {
	values := ""
	for i := 0; i < length; i++ {
		if i > 0 {
			values += ","
		}
		values += "(?, ?, ?, ?, ?, ?)"
	}
	return fmt.Sprintf(queryInsertTransactionDetail, values)
}

func NewStore(getDB func(dbName, replication string) (*sqlx.DB, error)) *Store {
	return &Store{
		getDB: getDB,
	}
}
