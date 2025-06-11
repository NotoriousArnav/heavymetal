APP_NAME := indexer
SRC := main.go
BUILD_DIR := build

.PHONY: all build run clean tidy

all: build

tidy:
	go mod tidy

build: tidy
	go build -o /tmp/$(APP_NAME) $(SRC)
	echo "Build complete: /tmp/$(APP_NAME)"

clean:
	rm -rf $(BUILD_DIR)

