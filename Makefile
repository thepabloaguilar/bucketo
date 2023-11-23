.PHONY: setup
setup:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

.PHONY: test
test:
	@go test -v -race -vet=all -count=1 -coverprofile=coverage.out ./...

.PHONY: lint
lint:
	@golangci-lint run ./...
