codegen:
	cd integration-svc && go generate ./...
	cd measurement-svc && go generate ./...

start-containers:
	docker compose -f compose.yaml up -d
	sleep 20

unit-tests:
	cd integration-svc && go run github.com/onsi/ginkgo/v2/ginkgo run --skip-package=e2e ./... && cd -
	cd measurement-svc && go run github.com/onsi/ginkgo/v2/ginkgo run --skip-package=e2e ./... && cd -

cleanup-containers:
	docker compose -f compose.yaml down -v
	docker rmi stellar-assignment-integration-svc:latest stellar-assignment-measurement-svc:latest stellar-assignment-api-gateway:latest

gateway-modbus-config:
	cp -f api-gateway/e2e/config/modbus_server.json config/modbus_server.json

gateway-e2e-tests: gateway-modbus-config start-containers
	sleep 10 && go run github.com/onsi/ginkgo/v2/ginkgo run api-gateway/e2e/...

measurement-modbus-config:
	cp -f measurement-svc/e2e/config/modbus_server.json config/modbus_server.json

measurement-e2e-tests: measurement-modbus-config start-containers
	sleep 10 && go run github.com/onsi/ginkgo/v2/ginkgo run measurement-svc/e2e/...

integration-modbus-config:
	cp -f integration-svc/e2e/config/modbus_server.json config/modbus_server.json

integration-e2e-tests: integration-modbus-config start-containers
	sleep 10 && go run github.com/onsi/ginkgo/v2/ginkgo run integration-svc/e2e/...