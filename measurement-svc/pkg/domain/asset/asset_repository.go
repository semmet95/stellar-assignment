package asset

import (
	"context"
	"fmt"

	"stellar-shared/pkg/domain/asset"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

// TODO: generate all the models using OAPI
type Asset struct {
	ID          string
	Setpoint    int64
	ActivePower int64
}

type AssetRepository interface {
	GetAssetByID(context.Context, string, string) (*Asset, error)
}

type assetRepository struct {
	queryAPI api.QueryAPI
}

// TODO: create a generic interface for the DB client
func NewAssetRepository(queryAPI api.QueryAPI) AssetRepository {
	return &assetRepository{
		queryAPI: queryAPI,
	}
}

// TODO: parameterize bucket name
func (ar *assetRepository) GetAssetByID(ctx context.Context, assetId, measurement string) (*Asset, error) {
	query := fmt.Sprintf(`from(bucket: "poc")
    |> range(start: -1s)
    |> filter(fn: (r) => r._measurement == "%s")
	|> filter(fn: (r) => r.id == "%s")
	|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
    |> last(column: "_time")`, measurement, assetId)

	result, err := ar.queryAPI.Query(ctx, query)
	if err == nil {
		if result.Err() != nil {
			return nil, fmt.Errorf("flux query error for asset id %s: %s", assetId, result.Err().Error())
		}

		if result.Next() { // get the single matched record
			record := result.Record()

			activePower, ok := record.ValueByKey(asset.ActivePowerKey).(int64) // Q: looks like influx only return int64?
			if !ok {
				return nil, fmt.Errorf("invalid active_power value: %v", activePower)
			}

			setpoint, ok := record.ValueByKey(asset.SetpointKey).(int64)
			if !ok {
				return nil, fmt.Errorf("invalid setpoint value: %v", setpoint)
			}

			return &Asset{
				ID:          assetId,
				Setpoint:    setpoint,
				ActivePower: activePower,
			}, nil
		}
		return nil, fmt.Errorf("no record matched for query: %s", query)
	}
	return nil, fmt.Errorf("query for asset id %s failed: %v", assetId, err)
}
