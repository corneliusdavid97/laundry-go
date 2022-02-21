package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/corneliusdavid97/laundry-go/src/customer"
	cust_handler "github.com/corneliusdavid97/laundry-go/src/customer/handler"
	cust_svc "github.com/corneliusdavid97/laundry-go/src/customer/service"
	cust_store "github.com/corneliusdavid97/laundry-go/src/customer/store"
	"github.com/corneliusdavid97/laundry-go/src/product"
	prod_handler "github.com/corneliusdavid97/laundry-go/src/product/handler"
	prod_svc "github.com/corneliusdavid97/laundry-go/src/product/service"
	prod_store "github.com/corneliusdavid97/laundry-go/src/product/store"
	"github.com/corneliusdavid97/laundry-go/src/transaction"
	trans_handler "github.com/corneliusdavid97/laundry-go/src/transaction/handler"
	trans_svc "github.com/corneliusdavid97/laundry-go/src/transaction/service"
	trans_store "github.com/corneliusdavid97/laundry-go/src/transaction/store"
	"github.com/corneliusdavid97/laundry-go/src/user"
	user_handler "github.com/corneliusdavid97/laundry-go/src/user/handler"
	user_svc "github.com/corneliusdavid97/laundry-go/src/user/service"
	user_store "github.com/corneliusdavid97/laundry-go/src/user/store"
	"github.com/corneliusdavid97/laundry-go/tools/postgresql"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	// context initialization
	timeout := time.Duration(5) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var err error

	// postgresql init
	err = postgresql.InitPostgresqlConfig(ctx, basepath)
	if err != nil {
		log.Fatalf("Failed to init postgresql database, err: %s", err.Error())
	}

	// user module
	{
		store := user_store.NewStore(func(dbName, replication string) (*sqlx.DB, error) {
			return postgresql.GetDB(dbName, replication)
		})
		svc := user_svc.NewService(store)
		user.Init(svc)
		userHTTPHandler := user_handler.NewHandler(svc, user_handler.Config{
			Timeout: time.Duration(3) * time.Second,
		})

		// handle HTTP request
		http.HandleFunc("/auth", userHTTPHandler.HandleAuthUser)
	}

	// customer module
	{
		store := cust_store.NewStore(func(dbName, replication string) (*sqlx.DB, error) {
			return postgresql.GetDB(dbName, replication)
		})
		svc := cust_svc.NewService(store)
		customer.Init(svc)
		userHTTPHandler := cust_handler.NewHandler(svc, cust_handler.Config{
			Timeout: time.Duration(3) * time.Second,
		})

		// handle HTTP request
		http.HandleFunc("/customer/all", userHTTPHandler.HandleGetAllActiveCustomer)
		http.HandleFunc("/customer/insert", userHTTPHandler.HandleInsertNewCustomer)
	}

	// product module
	{
		store := prod_store.NewStore(func(dbName, replication string) (*sqlx.DB, error) {
			return postgresql.GetDB(dbName, replication)
		})
		svc := prod_svc.NewService(store)
		product.Init(svc)
		userHTTPHandler := prod_handler.NewHandler(svc, prod_handler.Config{
			Timeout: time.Duration(3) * time.Second,
		})

		// handle HTTP request
		http.HandleFunc("/product/all", userHTTPHandler.HandleGetAllActiveProduct)
	}

	// transaction module
	{
		store := trans_store.NewStore(func(dbName, replication string) (*sqlx.DB, error) {
			return postgresql.GetDB(dbName, replication)
		})
		svc := trans_svc.NewService(store)
		transaction.Init(svc)
		userHTTPHandler := trans_handler.NewHandler(svc, trans_handler.Config{
			Timeout: time.Duration(3) * time.Second,
		})

		// handle HTTP request
		http.HandleFunc("/transaction/new", userHTTPHandler.HandleNewTransaction)
		http.HandleFunc("/transaction", userHTTPHandler.GetTransactionDataByID)
	}

	port := 4321
	log.Printf("Listening to port:%d", port)
	err = http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
