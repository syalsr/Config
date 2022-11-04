package server

import (
	"Config/proto"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

func StartGrpcServer(wg *sync.WaitGroup) {
	defer wg.Done()
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	proto.RegisterConfigWrapperServer(server, &ConfigWrapper{})
	log.Printf("gRPC server listening at %v", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
