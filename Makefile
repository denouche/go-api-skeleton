IS_LINUX=$(shell sed --version > /dev/null 2> /dev/null && echo $$?)
ifeq ($(IS_LINUX),0)
	SED_IN_PLACE=-i
else
	SED_IN_PLACE=-i ""
endif

NAME := $(shell cat info.yaml | sed -E '/title:/!d;s/.*: *'"'"'?([^$$])'"'"'?/\1/')
VERSION := $(shell cat info.yaml | sed -E '/version:/!d;s/.*: *'"'"'?([^$$])'"'"'?/\1/')
GITHASH := $(shell git rev-parse --short HEAD)
BUILDDATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

.PHONY: help
help: ## display this help
	@echo "Name: $(NAME)"
	@echo "Version: $(VERSION)"
	@echo "Git hash: $(GITHASH)"
	@echo "Build date: $(BUILDDATE)"
	@echo
	@echo "This is the list of available make targets:"
	@echo " $(shell cat Makefile | sed -E '/^[a-zA-Z-]+:.*##.*/!d;s/## *//;s/$$/\\n/')"

.PHONY: start
start: openapi ## start the application
	go run main.go --config config/local.json

.PHONY: start-offline
start-offline: openapi ## start the application in offline mode
	go run main.go --log-level debug --log-format text --db-in-memory

.PHONY: deps
deps: ## get the golang dependencies in the vendor folder
	GO111MODULE=on go mod vendor

.PHONY: build
build: ##  build the executable and set the version
	go build -o go-api-skeleton -ldflags "-X github.com/denouche/go-api-skeleton/handlers.ApplicationVersion=$(VERSION) -X github.com/denouche/go-api-skeleton/handlers.ApplicationName=$(NAME) -X github.com/denouche/go-api-skeleton/handlers.ApplicationGitHash=$(GITHASH) -X github.com/denouche/go-api-skeleton/handlers.ApplicationBuildDate=$(BUILDDATE)" main.go

.PHONY: test
test: ## run go test
	go test -v ./...

.PHONY: bump
bump: ## bump the version in the info.yaml file
	NEW_VERSION=`standard-version --dry-run | sed -E '/tagging release/!d;s/.*tagging release *v?(.*)/\1/g'`; \
		sed -E $(SED_IN_PLACE) 's/^(.*version: *).*$$/\1'$$NEW_VERSION'/' info.yaml

.PHONY: release
release: bump ## bump the version in the info.yaml, and make a release (commit, tag and push)
	git add info.yaml
	standard-version --message "chore(release): %s [ci skip]" --commit-all
	git push --follow-tags origin HEAD

.PHONY: openapi
openapi: ## install openapi-parser and generate openapi schema
	go get github.com/alexjomin/openapi-parser
	openapi-parser --output openapi.yaml
	openapi-parser merge --output openapi.yaml --main info.yaml --dir .
	go generate
