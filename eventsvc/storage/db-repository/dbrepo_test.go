package db_repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ladence/ecommerce-distributed-demo/eventstore"
	"reflect"
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

func TestGetEvents(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *eventstore.GetEventsRequest
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(dbmock sqlmock.Sqlmock)
		expected   []*eventstore.Event
		wantErr    bool
	}{
		{
			name: "Get Events (without filter)",
			args: args{
				context.TODO(),
				&eventstore.GetEventsRequest{},
			},
			beforeTest: func(dbmock sqlmock.Sqlmock) {
				dbmock.ExpectQuery("SELECT (.+)").WillReturnRows(sqlmock.NewRows([]string{"id", "eventtype", "aggregateid", "aggregatetype", "eventdata"}).AddRow("1", "order_created", "1", "order", "order_created_on_shelf"))
			},
			expected: []*eventstore.Event{
				&eventstore.Event{
					EventId:       "1",
					EventType:     "order_created",
					AggregateId:   "1",
					AggregateType: "order",
					EventData:     "order_created_on_shelf",
				},
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

			events, err := db.GetEvents(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbRepo.GetEvents wantErr: %v, err: %v", tt.wantErr, err)
			}
			if len(events) != len(tt.expected) {
				t.Errorf("dbRepo.GetEvents returned: %v events but expected: %v events", len(events), len(tt.expected))
			}
			for i := range events {
				if ok := reflect.DeepEqual(events[i], tt.expected[i]); !ok {
					t.Errorf("Returned event: %+v, but expected is: %+v", events[i], tt.expected[i])
				}
			}
		})
	}
}
