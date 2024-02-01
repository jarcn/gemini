BINARY_NAME=events-core

build:
	cd cmd/server; CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) -mod=vendor

dev:
	cd cmd/server; go build -o $(BINARY_NAME) -v
	cd cmd/server; ./$(BINARY_NAME) -config=./configs/config-dev.yaml