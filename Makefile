help:
	@echo "TODO"
	@echo "make start"

start:
	go run main.go --config config/local.json --environment development --port 8080

deps:
	GO111MODULE=on go mod vendor

build:
	go build -o go-api-skeleton main.go

