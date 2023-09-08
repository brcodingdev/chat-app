
APP_PKG = $(shell go list github.com/brcodingdev/service/internal/...)

lint:
	@echo "Linting"
	@golint -set_exit_status $(APP_PKG)
	@golangci-lint run --timeout 3m0s

test:
	@echo "Testing"
	@go test ./... -v -count=1 -race

build:
	@echo "Building docker image"
	@docker-compose build

run:
	@echo "Starting chat app"
	@docker-compose up -d
