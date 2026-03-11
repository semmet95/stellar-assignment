start-containers:
	docker compose -f compose.yaml up -d

cleanup-containers:
	docker compose -f compose.yaml down -v

start-stellar: start-containers
	go run cmd/app/main.go

acceptance-tests: start-stellar
	go run cmd/app/main.go