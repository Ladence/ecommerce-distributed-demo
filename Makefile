protoc:
	protoc eventstore/*.proto --go_out=. --go-grpc_out=.