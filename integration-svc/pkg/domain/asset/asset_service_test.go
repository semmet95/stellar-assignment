package asset_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"stellar-integration/pkg/domain/asset"
	"stellar-integration/pkg/domain/asset/assetfakes"
)

var _ = Describe("AssetService happy path", func() {
	It("calls repository to post asset", func() {
		fakeRepo := &assetfakes.FakeAssetRepository{}
		svc := asset.NewAssetService(fakeRepo)

		payload := &asset.Asset{ID: "test-id"}
		err := svc.PostAssetByID(context.Background(), payload, "measurement-sample")

		Expect(err).ToNot(HaveOccurred())
		Expect(fakeRepo.PostAssetByIDCallCount()).To(Equal(1))
		_, received, measurement := fakeRepo.PostAssetByIDArgsForCall(0)
		Expect(received).To(Equal(payload))
		Expect(measurement).To(Equal("measurement-sample"))
	})
})
