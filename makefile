APP_NAME=delivery

.PHONY: mod build test
build: test ## Build application
	mkdir -p build
	go build -o build/${APP_NAME} cmd/app/main.go
test: mod  ## Run tests
	go mockery
	go test ./...

mod:
	go mod tidy