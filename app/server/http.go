package server

import (
	"Config/app/proto"
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

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := proto.RegisterConfigWrapperHandlerFromEndpoint(ctx, gtw, os.Getenv("httpserver"), opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(":8081", gtw)
}

func StartHppServe(wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := grpc.Dial(
		os.Getenv("htppserver"),
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
