package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Ladence/ecommerce-distributed-demo/eventstore"
	"github.com/Ladence/ecommerce-distributed-demo/model"
	"github.com/Ladence/ecommerce-distributed-demo/ordersvc/storage"
	db_repository "github.com/Ladence/ecommerce-distributed-demo/ordersvc/storage/db-repository"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	subscribeSubj  = "ORDERS.paymentdebited"
	subscribeQueue = "order-validator"
	clientID       = "order-validator"
	eventSvcAddr   = "localhost:50001"
)

type PaymentDebitConsumer struct {
	ordersRepository storage.Repository
	ntCtx            nats.JetStreamContext
}

func (pdc *PaymentDebitConsumer) subscribeQueue(subject, queue string) error {
	_, err := pdc.ntCtx.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		msg.Ack()
		paymentDebitCmd := &model.PaymentDebitedCommand{}

		if err := json.Unmarshal(msg.Data, paymentDebitCmd); err != nil {
			log.Printf("error on umarshaling msg from queue order-validator: %v", err)
			return
		}
		changeOrderCmd := &model.ChangeOrderStatusCommand{
			OrderID: paymentDebitCmd.OrderID,
			Status:  "Approved",
		}

		// todo: change status via distributed repository
		log.Printf("changing order (%v) status to Approved", paymentDebitCmd.OrderID)
		time.Sleep(5 * time.Second)

		if err := executeOrderApprovedCommand(changeOrderCmd); err != nil {
			log.Printf("error occured on executing order approved command: %v", err)
			return
		}
	}, nats.Durable(clientID))
	return err
}

func executeOrderApprovedCommand(cmd *model.ChangeOrderStatusCommand) error {
	conn, err := grpc.Dial(eventSvcAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	esc := eventstore.NewEventSourceClient(conn)
	resp, err := esc.CreateEvent(context.Background(), &eventstore.CreateEventRequest{
		Event: &eventstore.Event{
			EventId:       uuid.NewString(),
			EventType:     "ORDERS.approved",
			AggregateId:   cmd.OrderID,
			AggregateType: "order",
			EventData:     "",
			Stream:        "ORDERS",
		},
	})
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("error on CreateEvent gRPC call: %v", resp.Error)
	}
	return nil
}

func newPaymentDebitConsumer(ntCtx nats.JetStreamContext, repository storage.Repository) *PaymentDebitConsumer {
	return &PaymentDebitConsumer{ntCtx: ntCtx, ordersRepository: repository}
}

func main() {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("error on connecting nats with default url: %v", err)
	}
	ntCtx, _ := conn.JetStream()

	cockroachRepo, err := db_repository.NewCockroachRepository("")
	if err != nil {
		log.Fatalf("error on creating cockroach repository: %v", err)
	}
	pdc := newPaymentDebitConsumer(ntCtx, cockroachRepo)
	if err := pdc.subscribeQueue(subscribeSubj, subscribeQueue); err != nil {
		log.Fatalf("error on subscribing queue: %v", err)
	}
}
