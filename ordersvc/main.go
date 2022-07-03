package main

import (
	"context"
	"encoding/json"
	"github.com/Ladence/ecommerce-distributed-demo/eventstore"
	"github.com/Ladence/ecommerce-distributed-demo/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	defaultAddr   = ":3001"
	apiVersion    = "v1"
	eventStoreSvc = "localhost:50001"
)

func createOrderEvent(order *model.Order, orderBytes []byte) *eventstore.Event {
	return &eventstore.Event{
		EventId:       uuid.NewString(),
		EventType:     "ORDERS.created",
		AggregateId:   order.ID,
		AggregateType: "order",
		EventData:     string(orderBytes),
		Stream:        "ORDERS",
	}
}

type orderHandler struct {
	esc eventstore.EventSourceClient
}

func newOrderHandler(esc eventstore.EventSourceClient) *orderHandler {
	return &orderHandler{esc: esc}
}

func (oh *orderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error on reading /order body, %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	order := &model.Order{}
	err = json.Unmarshal(bytes, order)
	if err != nil {
		log.Printf("error on unmarshalling /order body, %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := oh.esc.CreateEvent(context.Background(), &eventstore.CreateEventRequest{
		Event: createOrderEvent(order, bytes),
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			log.Printf("error from RPC server with: status code:%s message:%s", st.Code().String(), st.Message())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if !resp.Success {
		log.Printf("error on calling eventstore service (CreateEvent), %v", resp.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func initRoutes(oh *orderHandler) http.Handler {
	r := mux.NewRouter()
	r.Handle("api/"+apiVersion+"/order", oh).Methods(http.MethodPost)
	return r
}

func main() {
	conn, err := grpc.Dial(eventStoreSvc)
	if err != nil {
		log.Fatalf("error on establishing grpc connection: %v", err)
	}
	defer conn.Close()
	oh := newOrderHandler(eventstore.NewEventSourceClient(conn))
	s := http.Server{
		Addr:    defaultAddr,
		Handler: initRoutes(oh),
	}
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("error on starting server: %v", err)
	}
}
