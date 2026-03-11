package asset

import "context"

type AssetService interface {
	GetAssetByID(ctx context.Context, assetId, measurement string) (*Asset, error)
}

type assetService struct {
	assetRepo AssetRepository
}

func NewAssetService(assetRepo AssetRepository) AssetService {
	return &assetService{
		assetRepo: assetRepo,
	}
}

func (as *assetService) GetAssetByID(ctx context.Context, assetId, measurement string) (*Asset, error) {
	return as.assetRepo.GetAssetByID(ctx, assetId, measurement)
}
