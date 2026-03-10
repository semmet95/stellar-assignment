start-containers:
	docker compose -f compose.yaml up -d

cleanup-containers:
	docker compose -f compose.yaml down -v