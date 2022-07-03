package storage

import (
	"context"
	"github.com/Ladence/ecommerce-distributed-demo/model"
)

type Mock struct {
}

func (m Mock) CreateOrder(ctx context.Context, order model.Order) error {
	//TODO implement me
	panic("implement me")
}

func (m Mock) ChangeOrderStatus(ctx context.Context, cmd model.ChangeOrderStatusCommand) error {
	//TODO implement me
	panic("implement me")
}
