package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "stellar-measurement/gen"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	









	flag.Parse()

	conn, err := grpc.NewClient(
		*flag.String("addr", "localhost:50051", "the address to connect to"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewAssetServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := c.GetAsset(ctx, &pb.GetAssetRequest{Id: "871689260010377213"})
	if err != nil {
		log.Fatalf("could not get asset: %v", err)
	}
	log.Printf("active_power: %d, setpoint: %d", r.GetActivePower(), r.GetSetpoint())
}
