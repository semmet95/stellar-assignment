package listener

import (
	"context"
	"fmt"
	"log"
	"stellar-integration/pkg/domain/asset"
	"time"

	"github.com/simonvetter/modbus"
)

const (
	// Q: not sure how to make this dynamic
	assetID        = "871689260010377213"
	setpointKey    = "setpoint"
	activePowerKey = "active_power"
)

type assetListener struct {
	// TODO: add reference to asset svc
	modBusClient *modbus.ModbusClient
	pollInterval time.Duration
	assetSvc     asset.AssetService
}

type AssetListener interface {
	StartListening(ctx context.Context, unitID uint8) error
}

func NewAssetListener(mbClient *modbus.ModbusClient, pollInterval time.Duration, assetSvc asset.AssetService) AssetListener {
	return &assetListener{
		modBusClient: mbClient,
		pollInterval: pollInterval,
		assetSvc:     assetSvc,
	}
}

func (al *assetListener) StartListening(ctx context.Context, unitID uint8) error {
	// Q: How is this related to the asset ID?
	al.modBusClient.SetUnitId(unitID)

	err := al.modBusClient.Open()
	if err != nil {
		return fmt.Errorf("failed to open modbus connection: %v", err)
	}

	ticker := time.NewTicker(al.pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("modbus listener stopped")
			return nil
		case <-ticker.C:
			al.poll(ctx)
		}
	}
}

func (al *assetListener) poll(ctx context.Context) {
	// TODO: store register entries in a map and iterate over it
	// Q: ReadRegister will never return a negative value so how to validate here?
	setpoint, err := al.modBusClient.ReadRegister(30100, modbus.INPUT_REGISTER)
	if err != nil {
		log.Printf("failed to read setpoint from input register: %v\n", err)
		return
	}

	activePower, err := al.modBusClient.ReadRegister(40100, modbus.HOLDING_REGISTER)
	if err != nil {
		log.Printf("failed to read active_power from holding register: %v\n", err)
		return
	}

	al.handleRegisterValues(ctx, setpoint, activePower)
}

func (al *assetListener) handleRegisterValues(ctx context.Context, setpoint, activePower uint16) error {
	payload := &asset.Asset{
		Name:         "panel1",
		Type:         "SOLAR_PANEL",
		ID:           assetID,
		ConnProtocol: "TCP",
		RegisterMap: map[string]uint16{
			setpointKey:    setpoint,
			activePowerKey: activePower,
		},
	}
	return al.assetSvc.PostAssetByID(ctx, payload, assetID)
}
