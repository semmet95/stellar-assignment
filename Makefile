start-containers:
	docker compose -f compose.yaml up -d

cleanup-containers:
	docker compose -f compose.yaml down -v
	docker rmi stellar-assignment-integration-svc:latest stellar-assignment-measurement-svc:latest