BINARY_NAME=goto
BUILD_DIR=bin

.PHONY: build clean run

build-all:
	GOOS=linux   GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/$(BINARY_NAME).go
	GOOS=linux   GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/$(BINARY_NAME).go
	GOOS=darwin  GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/$(BINARY_NAME).go
	GOOS=darwin  GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/$(BINARY_NAME).go
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/$(BINARY_NAME).go

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/$(BINARY_NAME).go

clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
