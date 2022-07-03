package storage

import (
	"context"
	"github.com/Ladence/ecommerce-distributed-demo/model"
)

type Repository interface {
	CreateOrder(ctx context.Context, order model.Order) error
	ChangeOrderStatus(ctx context.Context, cmd model.ChangeOrderStatusCommand) error
}
