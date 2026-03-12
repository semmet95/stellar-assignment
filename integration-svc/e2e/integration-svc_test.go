package e2e

import (
	"context"
	"fmt"

	shared "stellar-shared/pkg/domain/asset"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	baseURL = "http://localhost:8086"
)

var (
	err error
)

var _ = Describe("Asset Integration API", func() {
	Context("Happy path", func() {
		When("Measurement service has read data from Modbus", func() {
			It("Influx DB has the correct measurements stored", func() {
				query := fmt.Sprintf(`from(bucket: "poc")
				|> range(start: -1s)
				|> filter(fn: (r) => r._measurement == "%s")
				|> filter(fn: (r) => r.id == "%s")
				|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")
				|> last(column: "_time")`, shared.Measurement, shared.AssetID)

				result, err := queryAPI.Query(context.TODO(), query)
				Expect(err).To(BeNil())
				Expect(result.Err()).To(BeNil())
				Expect(result.Next()).To(BeTrue())

				record := result.Record()
				activePower, ok := record.ValueByKey(shared.ActivePowerKey).(int64)
				Expect(ok).To(BeTrue())
				Expect(activePower).To(Equal(int64(17)))
				setpoint, ok := record.ValueByKey(shared.SetpointKey).(int64)
				Expect(ok).To(BeTrue())
				Expect(setpoint).To(Equal(int64(17)))
			})
		})
	})
})
