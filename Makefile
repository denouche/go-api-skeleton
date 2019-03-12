NAME := go-api-skeleton
VERSION := 0.0.1

.PHONY: help
help: ## display this help
	@echo "This is the list of available make targets:"
	@echo " $(shell cat Makefile | sed -r '/^[a-zA-Z-]+:.*##.*/!d;s/## *//;s/$$/\\n/')"

.PHONY: start
start: ## start the application
	go run main.go --config config/local.json

.PHONY: deps
deps: ## get the golang dependencies in the vendor folder
	GO111MODULE=on go mod vendor

.PHONY: build
build: ##  build the executable and set the version
	go build -o go-api-skeleton -ldflags "-X github.com/denouche/go-api-skeleton/handlers.ApplicationVersion=$(VERSION) -X github.com/denouche/go-api-skeleton/handlers.ApplicationName=$(NAME)" main.go

.PHONY: test
test: ## run go test
	go test -v ./...
