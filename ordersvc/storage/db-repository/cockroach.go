package db_repository

import (
	"context"
	"database/sql"

	"github.com/Ladence/ecommerce-distributed-demo/model"
	"github.com/Ladence/ecommerce-distributed-demo/ordersvc/storage"
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

func NewCockroachRepository(connectionString string) {

}

func setupSchema(db *sql.DB) error {
	return nil
}
