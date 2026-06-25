.PHONY: dev watch qa prod swagger tidy build install-swag install-air

DIST := dist/server

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

## One-shot run — compiles then runs binary so Ctrl+C exits cleanly (APP_ENV from .env)
dev: swagger
	@mkdir -p dist
	go build -o $(DIST) ./src/cmd/main.go
	-$(DIST)

## Override env to qa at runtime
qa: swagger
	@mkdir -p dist
	go build -o $(DIST) ./src/cmd/main.go
	-APP_ENV=qa $(DIST)

## Override env to prod at runtime
prod: swagger
	@mkdir -p dist
	go build -o $(DIST) ./src/cmd/main.go
	-APP_ENV=prod $(DIST)

## Build binary into dist/
build: swagger
	@mkdir -p dist
	go build -o $(DIST) ./src/cmd/main.go
