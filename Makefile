protoc:
	protoc eventstore/*.proto --go_out=. --go-grpc_out=.

build_eventsvc:
	go build -o .build/eventsvc eventsvc/server.go

build_ordersvc:
	go build -o .build/ordersvc ordersvc/main.go

build_paymentmanager:
	go build -o .build/paymentmanager paymentmanager/main.go

build_all: build_eventsvc build_paymentmanager build_ordersvc