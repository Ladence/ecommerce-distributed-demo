package db_repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Ladence/ecommerce-distributed-demo/model"
	"testing"
	"time"
)

func TestCreateOrder(t *testing.T) {
	type args struct {
		ctx   context.Context
		order model.Order
	}

	tests := []struct {
		name       string
		args       args
		beforeTest func(sqlmock sqlmock.Sqlmock)
		wantErr    bool
	}{
		{
			name: "create order",
			args: args{
				ctx: context.Background(),
				order: model.Order{ID: "1",
					CustomerID: "1",
					Status:     "Created",
					CreatedOn:  time.Date(2022, 07, 03, 18, 00, 00, 00, time.UTC),
					Amount:     1.0,
				},
			},
			beforeTest: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO (.+)").WithArgs("1", "1", "Created", "", "1.0").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		mockDB, mockSQL, _ := sqlmock.New()
		defer mockDB.Close()

		dbrepo := NewDbRepository(mockDB)
		if tt.beforeTest != nil {
			tt.beforeTest(mockSQL)
		}

		if err := dbrepo.CreateOrder(tt.args.ctx, tt.args.order); (err != nil) != tt.wantErr {
			t.Errorf("wantErr: %v but actual err is %v", tt.wantErr, err)
		}
	}
}
