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
	@go build -o ${BUILD_DIR}/${BINARY} ${GO_LINKER_FLAGS} main.go

deps:
	@printf "$(OK_COLOR)==> Downloading dependencies$(NO_COLOR)\n"
	@docker-compose exec slack-grand-race dep ensure

dev-up:
	@docker-compose up -d
	@docker-compose logs -f slack-grand-race

ssh:
	@docker-compose exec slack-grand-race /bin/sh

help:
	@echo "---------------------------------------------"
	@echo "List of available targets:"
	@echo "  build                      - Builds the binary and outputs it to out/folder    "
	@echo "  dev-up                     - Spins up the development containers"
	@echo "  deps                       - Downloads dependencies"
