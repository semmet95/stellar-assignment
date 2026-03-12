package e2e

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Measurement struct {
	ActivePower int64 `json:"activePower,string"`
	Setpoint    int64 `json:"setpoint,string"`
}

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Test Suite")
}
