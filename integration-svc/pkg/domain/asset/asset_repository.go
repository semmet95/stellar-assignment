package asset

import (
	"context"
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// TODO: generate all the models using OAPI
type Asset struct {
	Name string
	// TODO: custom types to represent a device type
	Type string
	ID   string
	// TODO: custom types to represent supported protocols
	ConnProtocol string
	// Q: Not sure about the map value type
	RegisterMap map[string]int16
}

//go:generate go tool counterfeiter . AssetRepository
type AssetRepository interface {
	PostAssetByID(ctx context.Context, asset *Asset, measurement string) error
}

type assetRepository struct {
	writer api.WriteAPI
}

// TODO: create a generic interface for the DB client
func NewAssetRepository(writer api.WriteAPI) AssetRepository {
	// Q: not sure about the best practice here
	go func() {
		for err := range writer.Errors() {
			log.Printf("influx writer error: %s\n", err.Error())
		}
	}()

	return &assetRepository{
		writer: writer,
	}
}

func (ar *assetRepository) PostAssetByID(ctx context.Context, asset *Asset, measurement string) error {
	point := influxdb2.NewPointWithMeasurement(measurement).
		AddTag("id", asset.ID).
		AddField("setpoint", asset.RegisterMap["setpoint"]).
		AddField("active_power", asset.RegisterMap["active_power"]).
		SetTime(time.Now())

	ar.writer.WritePoint(point)
	ar.writer.Flush()
	return nil
}
