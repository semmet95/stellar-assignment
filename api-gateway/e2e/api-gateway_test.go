package e2e

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	baseURL = "http://localhost:8081"
)

var (
	err error
)

var _ = Describe("Asset Measurements API", func() {
	Context("Testing GET endpoint", func() {
		When("GET request is sent with valid asset ID", func() {
			It("should return valid measurements", func() {
				client := &http.Client{}
				var measurement Measurement
				url := fmt.Sprintf("%s/asset/%s/measurements", baseURL, assetID)

				req, err := http.NewRequest("GET", url, nil)
				Expect(err).NotTo(HaveOccurred())

				req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml,application/json")
				resp, err := client.Do(req)
				Expect(err).NotTo(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				body, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				err = json.Unmarshal(body, &measurement)
				Expect(err).NotTo(HaveOccurred())
				Expect(measurement.ActivePower).To(BeEquivalentTo(1))
				Expect(measurement.Setpoint).To(BeEquivalentTo(1))
			})
		})
	})

	// TODO: add test for invalid asset ID after integration svc sends 404
	// TODO: add test to verify 5m caching
	// TODO: add test to verify the following flow: modbus register update -> cachd get response -> wait 5 mins -> updated get response
})
