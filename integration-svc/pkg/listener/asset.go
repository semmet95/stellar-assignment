package listener

import (
	"context"
	"fmt"
	"log"
	"stellar-integration/pkg/domain/asset"
	shared "stellar-shared/pkg/domain/asset"
	"time"

	"github.com/simonvetter/modbus"
)

type assetListener struct {
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

// StartListening opens the connection and polls.
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

// poll reads registers and forwards values.
func (al *assetListener) poll(ctx context.Context) {
	// TODO: store register entries in a map and iterate over it
	// TODO: add logging with different levels
	setpoint, err := al.modBusClient.ReadRegister(30100, modbus.INPUT_REGISTER)
	if err != nil {
		log.Printf("failed to read setpoint from input register: %v\n", err)
		return
	}
	signedSetpoint := int16(setpoint)
	if signedSetpoint < 0 {
		log.Printf("setpoint value is negative: %d\n", signedSetpoint)
		return
	}

	activePower, err := al.modBusClient.ReadRegister(40100, modbus.HOLDING_REGISTER)
	if err != nil {
		log.Printf("failed to read active_power from holding register: %v\n", err)
		return
	}
	signedactivePower := int16(activePower)
	if signedSetpoint < 0 {
		log.Printf("active_power value is negative: %d\n", signedactivePower)
		return
	}

	if signedactivePower > signedSetpoint {
		log.Printf("active_power value: %d is greater than setpoint: %d\n", signedactivePower, signedSetpoint)
		return
	}

	al.handleRegisterValues(ctx, signedSetpoint, signedactivePower)
}

// handleRegisterValues builds payload and posts it.
func (al *assetListener) handleRegisterValues(ctx context.Context, setpoint, activePower int16) error {
	payload := &asset.Asset{
		Name:         "panel1",
		Type:         "SOLAR_PANEL",
		ID:           shared.AssetID,
		ConnProtocol: "TCP",
		RegisterMap: map[string]int16{
			shared.SetpointKey:    setpoint,
			shared.ActivePowerKey: activePower,
		},
	}
	if err := al.assetSvc.PostAssetByID(ctx, payload, shared.Measurement); err != nil {
		log.Printf("failed to post measurements to service for asset %s: %v\n", payload.ID, err)
		return err
	}
	log.Printf("posted measurement for asset %s (active_power=%d setpoint=%d)\n", payload.ID, activePower, setpoint)
	return nil
}
