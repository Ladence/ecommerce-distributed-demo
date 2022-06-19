package storage

import (
	"context"
	"github.com/Ladence/ecommerce-distributed-demo/eventstore"
)

type Repository interface {
	CreateEvent(ctx context.Context, event *eventstore.Event) error
	GetEvents(ctx context.Context, filter *eventstore.GetEventsRequest) ([]*eventstore.Event, error)
}
