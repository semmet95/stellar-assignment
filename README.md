# Stellar Assignment

## Project overview
Small Go microservice suite that polls Modbus devices, writes/reads measurements in InfluxDB, exposes measurements over a gRPC gateway, and shares them via an HTTP → gRPC API gateway.

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
- Unit tests: `make unit-tests`. The project uses Ginkgo/Gomega for BDD-style tests in `*/pkg/...`.
- Fakes: generated fakes (via `counterfeiter`) live next to domain packages (e.g., `*/pkg/domain/asset/assetfakes`) and are used in unit tests.
- End-to-end: E2E suites are added for each service in the `e2e` folder. Each e2e suite has its own copy of Modbus server configuration, this way each suite tests against its own set of measurements.  
Run the e2e suites using the following make targets.
```
make integration-e2e-tests cleanup-containers
sleep 10
make measurement-e2e-tests cleanup-containers
sleep 10
make gateway-e2e-tests
```

## Local setup & deployment
Prerequisites: `go 1.26.x`, `docker`, `docker compose`, `make`.

Quick full-stack (compose):  
Compose file has default environment variable value set for each container. These values can be overridden as follows.
```
export INFLUX_HOST=localhost
export INFLUX_PORT=8086
export MODBUS_HOST=localhost
export MODBUS_PORT=5020
export MEASUREMENT_HOST=localhost
export MEASUREMENT_PORT=50051
export CACHE_DURATION_MINS=10
```
Run the following command to start all the services with the set configuration.
```bash
make start-containers
```
This uses [compose.yaml](compose.yaml) to start InfluxDB, a modbus-server, `integration-svc`, `measurement-svc`, and `api-gateway`.

Stop and clean:
```bash
make cleanup-containers
```
- Run services from source:
```bash
go run ./integration-svc/cmd/app
go run ./measurement-svc/cmd/app
go run ./api-gateway
```

## CI Workflows
- [Unit Tests](.github/workflows/unit-tests.yml): Runs all the unit tests when a PR to the `main` branch is created or when the `main` branch is updated.
- [E2E Tests](.github/workflows/e2e-tests.yml): Runs e2e test suites in all the services when a PR to the `main` branch is created or when the `main` branch is updated.