.PHONY: build
build:
	@go build -o ./dist/az-tui ./cmd/az-tui

.PHONY: test
test:
	@go test -race -cover ./...

.PHONY: lint
lint: ./bin/golangci-lint
	@./bin/golangci-lint run ./...

./bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v2.3.1
