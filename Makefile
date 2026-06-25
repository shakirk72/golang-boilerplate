.PHONY: dev watch qa prod swagger tidy build install-swag install-air

## Install swag CLI (one-time)
install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

## Install air live-reload (one-time)
install-air:
	go install github.com/air-verse/air@latest

## Generate / refresh docs from annotations
swagger:
	swag init -g src/cmd/main.go -o docs --quiet

## Download dependencies
tidy:
	go mod tidy

## Dev with live reload + auto swagger regen (reads APP_ENV from .env)
watch:
	air

## One-shot run — APP_ENV is read from .env (default: dev)
dev: swagger
	go run ./src/cmd/main.go

## Override env to qa at runtime
qa:
	APP_ENV=qa go run ./src/cmd/main.go

## Override env to prod at runtime
prod:
	APP_ENV=prod go run ./src/cmd/main.go

## Build binary
build: swagger
	go build -o bin/server ./src/cmd/main.go
