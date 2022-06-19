package storage

import (
	"context"
	"github.com/Ladence/ecommerce-distributed-demo/eventstore"
)

type Mock struct{}

func (Mock) CreateEvent(ctx context.Context, event *eventstore.Event) error {
	//TODO implement me
	panic("implement me")
}

func (Mock) GetEvents(ctx context.Context, filter *eventstore.GetEventsRequest) ([]*eventstore.Event, error) {
	//TODO implement me
	panic("implement me")
}
