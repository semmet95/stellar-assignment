package handler

import (
	"context"
	"fmt"

	pb "stellar-measurement/gen"
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

func (ah *AssetHandler) GetAsset(ctx context.Context, req *pb.GetAssetRequest) (*pb.AssetResponse, error) {
	assetId := req.GetId()
	asset, err := ah.assetSvc.GetAssetByID(ctx, assetId, shared.Measurement)
	if err != nil {
		return nil, fmt.Errorf("failed to get measurements for assed with id %s: %v", assetId, err)
	}

	return &pb.AssetResponse{
		ActivePower: int32(asset.ActivePower),
		Setpoint:    int32(asset.Setpoint),
	}, nil
}
