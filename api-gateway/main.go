package main

import (
	errcustom "api-gateway/pkg/error"
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

// TODO: need to filter error here otherwise it is added to the response
// main starts the HTTP->gRPC gateway.
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

	log.Printf("starting api gateway; forwarding to measurement service at %s:%s\n", measurementHost, measurementPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(errcustom.ErrorHandler),
		runtime.WithIncomingHeaderMatcher(headerMatcher),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Point the gateway at the gRPC server
	err := pb.RegisterAssetServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("%s:%s", measurementHost, measurementPort), opts)
	if err != nil {
		log.Fatalf("failed to register gateway to %s:%s: %v", measurementHost, measurementPort, err)
	}
	log.Println("gateway registered; starting HTTP server on :8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("api gateway stopped unexpectedly: %v", err)
	}
}

func headerMatcher(key string) (string, bool) {
	switch key {
	case "X-Client-Id":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
