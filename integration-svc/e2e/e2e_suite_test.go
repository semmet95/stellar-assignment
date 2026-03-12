package e2e

import (
	"testing"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	dbClient influxdb2.Client
	queryAPI api.QueryAPI
)

func TestE2E(t *testing.T) {
	var _ = BeforeSuite(func() {
		dbClient = influxdb2.NewClient("http://127.0.0.1:8086", "poc")
		queryAPI = dbClient.QueryAPI("poc")
	})

	var _ = AfterSuite(func() {
		dbClient.Close()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Test Suite")
}
