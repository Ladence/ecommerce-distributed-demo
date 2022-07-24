package db_repository

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Ladence/ecommerce-distributed-demo/model"
	"github.com/Ladence/ecommerce-distributed-demo/ordersvc/storage"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"
)

type CockroachRepository struct {
	db *sql.DB
}

func (c *CockroachRepository) CreateOrder(ctx context.Context, order model.Order) error {
	//TODO implement me
	panic("implement me")
}

func (c *CockroachRepository) ChangeOrderStatus(ctx context.Context, cmd model.ChangeOrderStatusCommand) error {
	//TODO implement me
	panic("implement me")
}

func NewDbRepository(db *sql.DB) storage.Repository {
	return &CockroachRepository{
		db: db,
	}
}

func NewCockroachRepository(connectionString string) (storage.Repository, error) {
	dbOrders, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	//if err = setupSchema(dbOrders); err != nil {
	//	return nil, err
	//}
	return &CockroachRepository{
		db: dbOrders,
	}, nil
}

func setupSchema(db *sql.DB) error {
	driver, err := cockroachdb.WithInstance(db, &cockroachdb.Config{})
	if err != nil {
		return err
	}
	fi, err := (&file.File{}).Open("./migrations/orders/")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("file", fi, "cockroachdb", driver)
	if err != nil {
		return err
	}
	return m.Steps(1)
}
