APP_NAME=app
API_MAIN=./cmd/api/main.go
WORKER_MAIN=./cmd/worker/main.go
MIGRATION_DIR=./db/migrations
SWAGGER_OUT=./docs/swagger

include .env

.PHONY: run-api
run-api:
	go run $(API_MAIN)

.PHONY: run-worker
run-worker:
	go run $(WORKER_MAIN)

.PHONY: build-api
build-api:
	go build -o bin/api $(API_MAIN)

.PHONY: build-worker
build-worker:
	go build -o bin/worker $(WORKER_MAIN)

.PHONY: test
test:
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: swagger
swagger:
	swag init -g ./cmd/api/main.go -o $(SWAGGER_OUT) --parseDependency --parseInternal

.PHONY: migrate-create
migrate-create:
	@read -p "migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $$name

.PHONY: migrate-up
migrate-up:
	@migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down:
	@migrate -path $(MIGRATION_DIR) -database "$(DATABASE_URL)" down 1

.PHONY: migrate-force
migrate-force:
	@read -p "version: " version; \
	migrate -path $(MIGRATION_DIR) -database "$$DATABASE_URL" force $$version
