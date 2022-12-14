package server

import (
	"Config/proto"
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func grpcGateway() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gtw := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := proto.RegisterConfigWrapperHandlerFromEndpoint(ctx, gtw, os.Getenv("grpcserver"), opts)
	if err != nil {
		return err
	}
	log.Println("Starting grpc gateway on :8081")
	
	return http.ListenAndServe(":8081", gtw)
}

func StartHppServe(wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := grpc.Dial(
		os.Getenv("grpcserver"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	if err = grpcGateway(); err != nil {
		log.Fatal(err)
	}

}
