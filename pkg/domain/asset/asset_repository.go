package asset

import (
	"context"
	"log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

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
		AddField("active_power", asset.RegisterMap["active_power"])

	ar.writer.WritePoint(point)
	ar.writer.Flush()
	return nil
}
