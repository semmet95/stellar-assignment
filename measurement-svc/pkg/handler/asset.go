package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "stellar-measurement/gen"
	"stellar-measurement/pkg/cache"
	"stellar-measurement/pkg/domain/asset"
	shared "stellar-shared/pkg/domain/asset"
)

type AssetHandler struct {
	pb.UnimplementedAssetServiceServer
	assetSvc asset.AssetService
}

func NewAssetHandler(svc asset.AssetService) *AssetHandler {
	return &AssetHandler{
		assetSvc: svc,
	}
}

// TODO: need to filter error here otherwise it is added to the response
func (ah *AssetHandler) GetAsset(ctx context.Context, req *pb.GetAssetRequest) (*pb.AssetResponse, error) {
	// check cache before connecting to DB
	clientId := ctx.Value("client_id").(string)
	cachedMeasurement := cache.GetCachedMeasurement(clientId)
	if cachedMeasurement != nil {
		log.Printf("cache hit for client %s\n", clientId)
		return &pb.AssetResponse{
			ActivePower: cachedMeasurement.ActivePower,
			Setpoint:    cachedMeasurement.Setpoint,
		}, nil
	}

	assetId := req.GetId()
	log.Printf("getting measurements for asset with id %s\n", assetId)
	asset, err := ah.assetSvc.GetAssetByID(ctx, assetId, shared.Measurement)
	if err != nil {
		return nil, fmt.Errorf("failed to get measurements for asset with id %s: %v", assetId, err)
	}

	// cache the response
	cache.UpdateCache(clientId, &cache.Measurement{
		ReqTime:     time.Now(),
		ActivePower: asset.ActivePower,
		Setpoint:    asset.Setpoint,
	})

	return &pb.AssetResponse{
		ActivePower: asset.ActivePower,
		Setpoint:    asset.Setpoint,
	}, nil
}
