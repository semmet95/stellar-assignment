package asset_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"stellar-measurement/pkg/domain/asset"
	"stellar-measurement/pkg/domain/asset/assetfakes"
)

var _ = Describe("AssetService happy path", func() {
	It("returns asset from repository", func() {
		fakeRepo := &assetfakes.FakeAssetRepository{}
		svc := asset.NewAssetService(fakeRepo)

		expected := &asset.Asset{ID: "a1", ActivePower: 100, Setpoint: 120}
		// configure fake to return expected
		fakeRepo.GetAssetByIDReturns(expected, nil)

		out, err := svc.GetAssetByID(context.Background(), "a1", "measurement")

		Expect(err).ToNot(HaveOccurred())
		Expect(out).To(Equal(expected))
		Expect(fakeRepo.GetAssetByIDCallCount()).To(Equal(1))
	})
})
