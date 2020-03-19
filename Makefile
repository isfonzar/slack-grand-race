### --------------------------------------------------------------------------------------------------------------------
### Variables
### (https://www.gnu.org/software/make/manual/html_node/Using-Variables.html#Using-Variables)
### --------------------------------------------------------------------------------------------------------------------

BUILD_DIR ?= out
BINARY=sgr

NAME=slack-grand-race
REPO=github.com/isfonzar/${NAME}

GO_LINKER_FLAGS=-ldflags="-s -w"
GO_PACKAGES := $(shell go list ./...)

# Other config
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

.PHONY: all test-unit test-integration

build:
	@printf "$(OK_COLOR)==> Building$(NO_COLOR)\n"
	@go build -o ${BUILD_DIR}/${BINARY} ${GO_LINKER_FLAGS} cmd/slack/main.go

deps:
	@printf "$(OK_COLOR)==> Downloading dependencies$(NO_COLOR)\n"
	@go mod vendor -v

dev-up:
	@docker-compose up -d
	@docker-compose logs -f slack-grand-race

up: deps
	@docker-compose up -d
	@docker-compose logs -f slack-grand-race

down:
	@docker-compose down

ssh:
	@docker-compose exec slack-grand-race /bin/sh

tests:
	@printf "$(OK_COLOR)==> Running tests$(NO_COLOR)\n"
	@go test -race ./...

tests-coverage:
	@printf "$(OK_COLOR)==> Running Unit tests$(NO_COLOR)\n"
	@go test -v -coverpkg=./... -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

help:
	@echo "---------------------------------------------"
	@echo "List of available targets:"
	@echo "  build                      - Builds the binary and outputs it to out/folder    "
	@echo "  dev-up                     - Spins up the development containers"
	@echo "  deps                       - Downloads dependencies"
	@echo "  tests                      - Executes unit tests"
	@echo "  tests-coverage             - Execute unit tests with test coverage"
