.PHONY: run test fmt vet up down

run:
	go run ./cmd/api

test:
	go test ./... -race -count=1

fmt:
	go fmt ./...

vet:
	go vet ./...

up:
	docker compose -f deployments/docker-compose.yml up -d

down:
	docker compose -f deployments/docker-compose.yml down
