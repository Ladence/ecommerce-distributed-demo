package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Ladence/ecommerce-distributed-demo/eventstore"
	"github.com/Ladence/ecommerce-distributed-demo/model"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"log"
)

const (
	eventStoreSvcUri = "localhost:50001"

	consumerName = "payment-manager"
)

type orderEventConsumer struct {
	ntCtx nats.JetStreamContext
}

func newOrderEventConsumer(ntCtx nats.JetStreamContext) *orderEventConsumer {
	return &orderEventConsumer{ntCtx: ntCtx}
}

func (c *orderEventConsumer) subscribeQueue(subj, queueName string) error {
	_, err := c.ntCtx.QueueSubscribe(subj, queueName, func(msg *nats.Msg) {
		order := &model.Order{}
		if err := json.Unmarshal(msg.Data, order); err != nil {
			return
		}
		log.Printf("message order_created received: %v", order)

		cmd := model.PaymentDebitedCommand{
			OrderID:    order.ID,
			CustomerID: order.CustomerID,
			Amount:     order.Amount,
		}
		if err := executePaymentDebitedCommand(&cmd); err != nil {
			log.Printf("error on executing payment debited command. err: %v", err)
			return
		}
	}, nats.Durable(consumerName), nats.ManualAck())
	return err
}

func executePaymentDebitedCommand(cmd *model.PaymentDebitedCommand) error {
	conn, err := grpc.Dial(eventStoreSvcUri, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	paymentData, _ := json.Marshal(cmd)
	esclient := eventstore.NewEventSourceClient(conn)
	event := eventstore.Event{
		EventId:       uuid.NewString(),
		EventType:     "ORDERS.paymentdebited",
		AggregateId:   cmd.OrderID,
		AggregateType: "order",
		EventData:     string(paymentData),
		Stream:        "ORDERS",
	}
	resp, err := esclient.CreateEvent(context.Background(), &eventstore.CreateEventRequest{
		Event: &event,
	})
	if err != nil {
		return err
	}
	if resp.Success {
		return nil
	}
	return fmt.Errorf("err from grpc event service: %s", resp.Error)
}

func main() {
	conn, err := nats.Connect("0.0.0.0:4222")
	if err != nil {
		log.Fatalf("Error on establishing connection to nats")
	}
	ntCtx, _ := conn.JetStream()

	eventConsumer := newOrderEventConsumer(ntCtx)
	if err := eventConsumer.subscribeQueue("ORDERS.created", "payment-manager"); err != nil {
		log.Fatalf("error on subscribeQueue: %v", err)
	}
}
