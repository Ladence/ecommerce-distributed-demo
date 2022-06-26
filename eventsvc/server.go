package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Ladence/ecommerce-distributed-demo/eventstore"
	"github.com/Ladence/ecommerce-distributed-demo/eventsvc/storage"
	db_repository "github.com/Ladence/ecommerce-distributed-demo/eventsvc/storage/db-repository"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type Server struct {
	eventstore.UnimplementedEventSourceServer
	store   storage.Repository
	natsCtx nats.JetStreamContext
}

func NewServer(store storage.Repository, streamContext nats.JetStreamContext) *Server {
	return &Server{store: store, natsCtx: streamContext}
}

func (s *Server) CreateEvent(ctx context.Context, request *eventstore.CreateEventRequest) (*eventstore.CreateEventResponse, error) {
	if err := s.store.CreateEvent(ctx, request.Event); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	go publishEvent(s.natsCtx, request.Event)
	return &eventstore.CreateEventResponse{Success: true}, nil
}

func (s *Server) GetEvents(ctx context.Context, request *eventstore.GetEventsRequest) (*eventstore.GetEventsResponse, error) {
	events, err := s.store.GetEvents(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &eventstore.GetEventsResponse{Events: events}, nil
}

func publishEvent(natsCtx nats.JetStreamContext, event *eventstore.Event) {
	ack, err := natsCtx.Publish(event.EventType, []byte(event.EventData))
	if err != nil {
		log.Printf("error on publishing event (+%v), err: %v\n", event, err)
		return
	}
	log.Printf("succesfully published an event! stream: %s\n", ack.Stream)
}

func main() {
	port := flag.String("port", "50001", "port for grpc server")
	flag.Parse()

	conn, _ := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	grpcServer := grpc.NewServer()

	ntConn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("error on connecting natsCtx: %v", err)
	}

	ntCtx, err := ntConn.JetStream()
	if err != nil {
		log.Fatalf("error on perceiving natsCtx JetStream context: %v", err)
	}

	repo, err := db_repository.NewSqliteRepository("events.db")
	if err != nil {
		log.Fatalf("error on creating sqlite repository: %v", err)
	}

	server := NewServer(repo, ntCtx)
	eventstore.RegisterEventSourceServer(grpcServer, server)
	log.Println("registered a grpc server")
	if err := grpcServer.Serve(conn); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
