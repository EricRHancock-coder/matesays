BINARY_NAME=matesays
BUILD_DIR=bin
MAIN_PATH=./matesays.go

.PHONY: all run fmt build clean

all: fmt build run

run:
	go run $(MAIN_PATH)

fmt:
	go fmt $(MAIN_PATH)

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH) 

clean:
	rm -rf $(BUILD_DIR)
