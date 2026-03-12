start-containers:
	docker compose -f compose.yaml up -d

cleanup-containers:
	docker compose -f compose.yaml down -v
	docker rmi stellar-assignment-integration-svc:latest stellar-assignment-measurement-svc:latest stellar-assignment-api-gateway:latest

e2e-tests: start-containers
	sleep 10 && go run github.com/onsi/ginkgo/v2/ginkgo run api-gateway/e2e/...