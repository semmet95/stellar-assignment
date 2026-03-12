package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"stellar-integration/pkg/domain/asset"
	"stellar-integration/pkg/listener"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/simonvetter/modbus"
)

const (
	// Q: not sure how to configure it dynamically
	unitID = 1
)

var (
	dbClient influxdb2.Client
	mbClient *modbus.ModbusClient
	err      error
)

func main() {
	// TODO: move config setup logic to dedicated conf package stored in the shared module
	// initialize config
	modbusHost, ok := os.LookupEnv("MODBUS_HOST")
	if !ok {
		log.Fatal("MODBUS_HOST environment variable not set")
	}

	modbusPort, ok := os.LookupEnv("MODBUS_PORT")
	if !ok {
		log.Fatal("MODBUS_PORT environment variable not set")
	}

	influxHost, ok := os.LookupEnv("INFLUX_HOST")
	if !ok {
		log.Fatal("MODBUS_HOST environment variable not set")
	}

	influxPort, ok := os.LookupEnv("INFLUX_PORT")
	if !ok {
		log.Fatal("MODBUS_PORT environment variable not set")
	}

	log.Printf("connecting to Modbus at %s:%s\n", modbusHost, modbusPort)
	if err := initModbusClient(modbusHost, modbusPort); err != nil {
		log.Fatalf("failed to initialize modbus client at %s:%s: %v", modbusHost, modbusPort, err)
	}
	defer mbClient.Close()

	log.Printf("connecting to InfluxDB at %s:%s\n", influxHost, influxPort)
	dbClient = influxdb2.NewClient(fmt.Sprintf("http://%s:%s", influxHost, influxPort), "poc")
	defer dbClient.Close()

	assetRepo := asset.NewAssetRepository(dbClient.WriteAPI("poc", "poc"))
	assetSvc := asset.NewAssetService(assetRepo)

	// start the listener
	log.Println("starting modbus listener")
	if err := listener.NewAssetListener(
		mbClient,
		1*time.Second,
		assetSvc,
	).StartListening(context.Background(), unitID); err != nil {
		log.Fatalf("modbus listener failed to start with error: %v", err)
	}
}

func initModbusClient(host, port string) error {
	mbClient, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     fmt.Sprintf("tcp://%s:%s", host, port),
		Timeout: 500 * time.Millisecond,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize modbus client: %v", err)
	}

	if err = mbClient.Open(); err != nil {
		return fmt.Errorf("failed to connect to modbus: %v", err)
	}

	return nil
}
