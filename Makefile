include .env

.PHONY: *

build:
	go build -a -installsuffix cgo -o app ./cmd/main/main.go

test:
	go test ./...

lint:
	golangci-lint run ./...

up-env-local:
	docker compose -f ./docker/docker-compose.dev.yml --env-file=.env up -d

down-env-local:
	docker compose -f ./docker/docker-compose.dev.yml --env-file=.env stop

migrate:
	docker run --rm -it --network=${DC_NETWORK} -v "$(PWD)/migrations:/db/migrations" -e DATABASE_URL=$(MIGRATE_DATABASE_URL) --env-file=.env ghcr.io/amacneil/dbmate $(command)
