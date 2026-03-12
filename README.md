# Stellar Assignment

## Project overview
Small Go microservice suite that polls Modbus devices, writes/reads measurements in InfluxDB, and exposes measurements over an HTTP/gRPC gateway.

## Modules
- **api-gateway:** [api-gateway](api-gateway) — HTTP → gRPC gateway (grpc-gateway) that forwards REST calls to the measurement gRPC service.
- **integration-svc:** [integration-svc](integration-svc) — Modbus poller and writer; polls devices and writes points to InfluxDB.
- **measurement-svc:** [measurement-svc](measurement-svc) — gRPC measurement service that queries InfluxDB and serves measurement data.
- **shared:** [shared](shared) — shared domain types/constants used across modules.
- **Repo-level orchestration:** `compose.yaml`, `Dockerfile.gateway`, `Dockerfile.integration`, `Dockerfile.measurement` and [Makefile](Makefile) for local orchestration.

## How modules connect
- `integration-svc` polls Modbus devices and writes measurement points into InfluxDB.
- `measurement-svc` queries InfluxDB for the latest measurement and exposes it via gRPC.
- `api-gateway` forwards HTTP requests to `measurement-svc` using a gRPC-gateway.
- `compose.yaml` wires InfluxDB, a modbus-server, and the services for local end-to-end runs.

## Configuration / env vars
- `MODBUS_HOST`, `MODBUS_PORT` — for `integration-svc` to reach the Modbus server.
- `INFLUX_HOST`, `INFLUX_PORT` — InfluxDB address used by both services.
- `MEASUREMENT_HOST`, `MEASUREMENT_PORT` — used by `api-gateway` to reach `measurement-svc`.
- `CACHE_DURATION_MINS` — cache TTL used by `measurement-svc`.

Set these in your shell or in your container runtime before running services locally.

## Testing
- Unit tests: `go test ./...`. The project uses Ginkgo/Gomega for BDD-style tests in `*/pkg/...`.
- Fakes: generated fakes (via `counterfeiter`) live next to domain packages (e.g., `*/pkg/domain/asset/assetfakes`) and are used in unit tests.
- End-to-end: E2E suites are under [api-gateway/e2e](api-gateway/e2e). Run `make e2e-tests` to start containers and execute the e2e suites.
- Suggested GitHub workflows (CI):
  - `unit` job: `go test ./...`, `go vet`/linters, build binaries.
  - `e2e` job (optional/separate): build Docker images and run `make e2e-tests` against a composed environment.

## Local setup & deployment
Prerequisites: `go 1.26.x`, `docker`, `docker compose`, `make`.

Quick full-stack (compose):
```bash
make start-containers
```
This uses [compose.yaml](compose.yaml) to start InfluxDB, a modbus-server, `integration-svc`, `measurement-svc`, and `api-gateway`.

Stop and clean:
```bash
make cleanup-containers
```

Run services locally (development):
- Start InfluxDB and Modbus server via compose or external services.
- Export env vars (example):
```bash
export INFLUX_HOST=localhost
export INFLUX_PORT=8086
export MODBUS_HOST=localhost
export MODBUS_PORT=5020
export MEASUREMENT_HOST=localhost
export MEASUREMENT_PORT=50051
```
- Run services from source:
```bash
go run ./integration-svc/cmd/app
go run ./measurement-svc/cmd/app
go run ./api-gateway
```

Run tests:
- Unit tests:
```bash
go test ./...
```
- E2E (starts containers then runs suites):
```bash
make e2e-tests
```

## Notes & tips
- Logging: services use the standard `log` package for concise startup, connection and error context logs.
- Regenerating fakes: use `counterfeiter` (referenced in module tooling) to update fakes.
- CI: keep unit and e2e separate (e2e usually runs in a dedicated job with service containers).

---
Concise README created to help new contributors run and test the project locally. If you want, I can also add a GitHub Actions workflow example or expand any section further.