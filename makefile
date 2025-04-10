APP_NAME=nats-prom-bridge
SRC=./cmd/exporter
BIN_DIR=./bin

# Platform-specific file names
WIN_BIN=$(BIN_DIR)/$(APP_NAME).exe
LINUX_BIN=$(BIN_DIR)/$(APP_NAME)-linux
MAC_BIN=$(BIN_DIR)/$(APP_NAME)-mac

.PHONY: all build windows linux mac clean

all: windows linux mac

build:
	go build -o $(BIN_DIR)/$(APP_NAME) $(SRC)

windows:
	GOOS=windows GOARCH=amd64 go build -o $(WIN_BIN) $(SRC)

linux:
	GOOS=linux GOARCH=amd64 go build -o $(LINUX_BIN) $(SRC)

mac:
	GOOS=darwin GOARCH=amd64 go build -o $(MAC_BIN) $(SRC)

clean:
	rm -rf $(BIN_DIR)
