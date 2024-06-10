$(shell cp -n .env.example .env)
include .env
export

run:
	@go run main.go
.PHONY: run

test:
	@go test ./...
.PHONY: test
