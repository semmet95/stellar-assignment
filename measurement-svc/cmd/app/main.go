package main

import (
	"fmt"
	"log"
	"net"
	"os"
	pb "stellar-measurement/gen"
	"stellar-measurement/pkg/domain/asset"
	"stellar-measurement/pkg/handler"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"google.golang.org/grpc"
)

var (
	dbClient influxdb2.Client
	err      error
)

// TODO: need to find a way to generate protobufs dynamically and share proto files
func main() {

	influxHost, ok := os.LookupEnv("INFLUX_HOST")
	if !ok {
		log.Fatal("MODBUS_HOST environment variable not set")
	}

	influxPort, ok := os.LookupEnv("INFLUX_PORT")
	if !ok {
		log.Fatal("MODBUS_PORT environment variable not set")
	}

	dbClient = influxdb2.NewClient(fmt.Sprintf("http://%s:%s", influxHost, influxPort), "poc")
	defer dbClient.Close()

	assetRepo := asset.NewAssetRepository(dbClient.QueryAPI("poc"))
	assetSvc := asset.NewAssetService(assetRepo)
	assetHandler := handler.NewAssetHandler(assetSvc)

	server := grpc.NewServer()
	pb.RegisterAssetServiceServer(server, assetHandler)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to start tcp listener: %v", err)
	}

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to start grps server: %v", err)
	}
}
