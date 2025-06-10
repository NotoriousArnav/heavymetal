APP_NAME := indexer
SRC := main.go
BUILD_DIR := build

.PHONY: all build run clean tidy

all: build

tidy:
	go mod tidy

build: tidy
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC)

clean:
	rm -rf $(BUILD_DIR)

