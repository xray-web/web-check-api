run:
	@go run main.go
.PHONY: run

test:
	@go test ./...
.PHONY: test