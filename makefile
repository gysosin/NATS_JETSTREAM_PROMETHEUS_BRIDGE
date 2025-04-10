APP_NAME=nats-prom-bridge
SRC=./cmd/exporter
BIN_DIR=./bin

WIN_BIN=$(BIN_DIR)/$(APP_NAME).exe
LINUX_BIN=$(BIN_DIR)/$(APP_NAME)-linux
MAC_BIN=$(BIN_DIR)/$(APP_NAME)-mac

.PHONY: all build windows linux mac clean

all: windows linux mac

# Auto-detect OS to build for current system
build:
ifeq ($(OS),Windows_NT)
	go build -o $(WIN_BIN) $(SRC)
else
	go build -o $(BIN_DIR)/$(APP_NAME) $(SRC)
endif

# Use platform-specific env vars for cross-compilation
windows:
	@echo "Building Windows binary..."
	go env -w GOOS=windows GOARCH=amd64
	go build -o $(WIN_BIN) $(SRC)
	go env -u GOOS GOARCH

linux:
	@echo "Building Linux binary..."
	go env -w GOOS=linux GOARCH=amd64
	go build -o $(LINUX_BIN) $(SRC)
	go env -u GOOS GOARCH

mac:
	@echo "Building macOS binary..."
	go env -w GOOS=darwin GOARCH=amd64
	go build -o $(MAC_BIN) $(SRC)
	go env -u GOOS GOARCH

clean:
	@echo "Cleaning binaries..."
	@if exist $(BIN_DIR) rmdir /s /q $(BIN_DIR) || rm -rf $(BIN_DIR)
