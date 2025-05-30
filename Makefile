BINARY_NAME=goto
BUILD_DIR=bin

.PHONY: build clean run

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/$(BINARY_NAME).go

run: build
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
