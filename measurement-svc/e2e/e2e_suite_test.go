package e2e

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "stellar-measurement/gen"
)

var (
	client pb.AssetServiceClient
	conn   *grpc.ClientConn
	err    error
)

type Measurement struct {
	ActivePower int64 `json:"activePower,string"`
	Setpoint    int64 `json:"setpoint,string"`
}

func TestE2E(t *testing.T) {
	var _ = BeforeSuite(func() {
		conn, err = grpc.NewClient(
			"localhost:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).To(BeNil())
		client = pb.NewAssetServiceClient(conn)
	})

	var _ = AfterSuite(func() {
		conn.Close()
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Test Suite")
}
