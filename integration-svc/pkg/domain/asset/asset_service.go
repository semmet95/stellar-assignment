package asset

import "context"

type AssetService interface {
	PostAssetByID(ctx context.Context, asset *Asset, measurement string) error
}

type assetService struct {
	assetRepo AssetRepository
}

func NewAssetService(assetRepo AssetRepository) AssetService {
	return &assetService{
		assetRepo: assetRepo,
	}
}

func (as *assetService) PostAssetByID(ctx context.Context, asset *Asset, measurement string) error {
	as.assetRepo.PostAssetByID(ctx, asset, measurement)

	return nil
}
