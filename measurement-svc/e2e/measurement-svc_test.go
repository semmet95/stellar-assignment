package e2e

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	pb "stellar-measurement/gen"
	shared "stellar-shared/pkg/domain/asset"
)

var _ = Describe("Asset Measurement API", func() {
	Context("Happy path", func() {
		When("Integration service has written data to Influx DB", func() {
			It("Measuremet service can share that data via a gRPC endpoint", func() {
				ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
				defer cancel()

				assetMeasurement, err := client.GetAsset(ctx, &pb.GetAssetRequest{Id: shared.AssetID})
				Expect(err).To(BeNil())

				Expect(assetMeasurement.ActivePower).To(BeEquivalentTo(int64(4369)))
				Expect(assetMeasurement.Setpoint).To(BeEquivalentTo(int64(4420)))
			})
		})
	})
})
