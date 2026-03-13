package cache_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cache "stellar-measurement/pkg/cache"
)

var _ = Describe("Cache", func() {
	It("stores and retrieves measurement (happy path)", func() {
		now := time.Now()
		m := &cache.Measurement{
			ReqTime:     now,
			ActivePower: 123,
			Setpoint:    456,
		}

		clientID := "client1"
		assetID := "assetA"

		cache.UpdateCache(clientID, assetID, m)

		got := cache.GetCachedMeasurement(clientID, assetID)
		Expect(got).ToNot(BeNil())
		Expect(got.ActivePower).To(Equal(m.ActivePower))
		Expect(got.Setpoint).To(Equal(m.Setpoint))
		// ensure cache is considered fresh (uses package default)
		Expect(time.Since(got.ReqTime)).To(BeNumerically("<=", 5*time.Minute))
	})

	It("expires measurements older than cache duration", func() {
		old := time.Now().Add(-10 * time.Minute)
		m := &cache.Measurement{
			ReqTime:     old,
			ActivePower: 1,
			Setpoint:    2,
		}

		clientID := "client-expire"
		assetID := "asset-expire"

		cache.UpdateCache(clientID, assetID, m)

		got := cache.GetCachedMeasurement(clientID, assetID)
		Expect(got).To(BeNil())
	})

	It("should cache for unique combination of client and asset IDs", func() {
		now := time.Now()
		m := &cache.Measurement{
			ReqTime:     now,
			ActivePower: 123,
			Setpoint:    456,
		}

		clientID := "client1"
		assetID := "assetA"

		cache.UpdateCache(clientID, assetID, m)

		got := cache.GetCachedMeasurement(clientID, "assetB")
		Expect(got).To(BeNil())
	})
})
