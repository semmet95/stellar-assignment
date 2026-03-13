package asset

import "context"

type AssetService interface {
	PostAssetByID(ctx context.Context, asset *Asset, measurement string) error
}

type assetService struct {
	assetRepo AssetRepository
}

// NewAssetService creates an AssetService.
func NewAssetService(assetRepo AssetRepository) AssetService {
	return &assetService{
		assetRepo: assetRepo,
	}
}

// PostAssetByID posts asset measurement via repository.
func (as *assetService) PostAssetByID(ctx context.Context, asset *Asset, measurement string) error {
	return as.assetRepo.PostAssetByID(ctx, asset, measurement)
}
