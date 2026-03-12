package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	pb "stellar-measurement/gen"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// TODO: move config setup logic to dedicated conf package stored in the shared module
	// initialize config
	measurementHost, ok := os.LookupEnv("MEASUREMENT_HOST")
	if !ok {
		log.Fatal("MEASUREMENT_HOST environment variable not set")
	}

	measurementPort, ok := os.LookupEnv("MEASUREMENT_PORT")
	if !ok {
		log.Fatal("MEASUREMENT_PORT environment variable not set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Point the gateway at the gRPC server
	err := pb.RegisterAssetServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%s", measurementHost, measurementPort), opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("failed to start gateway on port 8080 : %v", err)
	}
}
