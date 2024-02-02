BINARY_NAME=gemini

build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)

dev:
	go build -o $(BINARY_NAME) -v
	#./$(BINARY_NAME) -config=./configs/config-dev.yaml