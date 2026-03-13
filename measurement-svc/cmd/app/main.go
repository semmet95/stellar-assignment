package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	pb "stellar-measurement/gen"
	"stellar-measurement/pkg/domain/asset"
	"stellar-measurement/pkg/handler"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	dbClient influxdb2.Client
	err      error
)

// TODO: need to find a way to generate protobufs dynamically and share proto files
// main starts the measurement gRPC server.
func main() {

	influxHost, ok := os.LookupEnv("INFLUX_HOST")
	if !ok {
		log.Fatal("INFLUX_HOST environment variable not set")
	}

	influxPort, ok := os.LookupEnv("INFLUX_PORT")
	if !ok {
		log.Fatal("INFLUX_PORT environment variable not set")
	}

	log.Printf("connecting to InfluxDB at %s:%s\n", influxHost, influxPort)
	dbClient = influxdb2.NewClient(fmt.Sprintf("http://%s:%s", influxHost, influxPort), "poc")
	defer dbClient.Close()

	assetRepo := asset.NewAssetRepository(dbClient.QueryAPI("poc"))
	assetSvc := asset.NewAssetService(assetRepo)
	assetHandler := handler.NewAssetHandler(assetSvc)

	server := grpc.NewServer(grpc.UnaryInterceptor(clientIDInterceptor))
	pb.RegisterAssetServiceServer(server, assetHandler)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to start tcp listener on :50051: %v", err)
	}

	log.Println("measurement service listening on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("measurement service stopped unexpectedly: %v", err)
	}
}

// clientIDInterceptor forwards headers from incoming http request to grpc request metadata
func clientIDInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	clientID := "unknown"

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get("x-client-id"); len(values) > 0 {
			clientID = values[0]
		}
	}

	ctx = context.WithValue(ctx, "client_id", clientID)

	return handler(ctx, req)
}
