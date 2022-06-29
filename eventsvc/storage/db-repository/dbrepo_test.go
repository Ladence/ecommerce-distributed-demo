package db_repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ladence/ecommerce-distributed-demo/eventstore"
	"testing"
)

func TestCreateEvent(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *eventstore.Event
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "success create input (without event data)",
			args: args{
				ctx: context.TODO(),
				input: &eventstore.Event{
					EventId:       "1",
					EventType:     "order_created",
					AggregateId:   "1",
					AggregateType: "order",
					Stream:        "todo",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectExec("INSERT INTO (.+)").WithArgs("1", "order_created", "1", "order", "todo").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "success create input (with event data)",
			args: args{
				ctx: context.TODO(),
				input: &eventstore.Event{
					EventId:       "1",
					EventType:     "order_created",
					EventData:     "publishing buy-order",
					AggregateId:   "1",
					AggregateType: "order",
					Stream:        "todo",
				},
			},
			beforeTest: func(s sqlmock.Sqlmock) {
				s.ExpectExec("INSERT (.*)").WithArgs("1", "order_created", "1", "order", "publishing buy-order", "todo").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, _ := sqlmock.New()
			defer mockDB.Close()

			db := NewDbRepository(mockDB)

			if tt.beforeTest != nil {
				tt.beforeTest(mockSQL)
			}

			err := db.CreateEvent(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbRepo.CreateEvent wantErr: %v, err: %v", tt.wantErr, err)
			}
		})
	}
}
