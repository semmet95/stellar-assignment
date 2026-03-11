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
	// TODO: move config setup logic to dedicated conf package with graceful exit
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

	if err := initModbusClient(modbusHost, modbusPort); err != nil {
		log.Fatal(err)
	}
	defer mbClient.Close()

	dbClient = influxdb2.NewClient(fmt.Sprintf("http://%s:%s", influxHost, influxPort), "poc")
	defer dbClient.Close()

	assetRepo := asset.NewAssetRepository(dbClient.WriteAPI("poc", "poc"))
	assetSvc := asset.NewAssetService(assetRepo)

	// start the listener
	listener.NewAssetListener(
		mbClient,
		1*time.Second,
		assetSvc,
	).StartListening(context.Background(), unitID) // could start this in a goroutine if a server needs to run
}

func initModbusClient(host, port string) error {
	mbClient, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     fmt.Sprintf("tcp://%s:%s", host, port),
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize modbus client: %v", err)
	}

	if err = mbClient.Open(); err != nil {
		return fmt.Errorf("failed to connect to modbus: %v", err)
	}

	return nil
}
