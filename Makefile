.PHONY: build clean test fmt lint

BINARY_NAME=scion
BUILD_DIR=build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/scion

clean:
	rm -rf $(BUILD_DIR)

test:
	go test ./...

test-v:
	go test -v ./...

test-cover:
	go test -cover ./...

fmt:
	go fmt ./...

lint:
	go vet ./...
