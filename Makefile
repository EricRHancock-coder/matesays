BINARY_NAME=matesays
BUILD_DIR=bin
MAIN_PATH=./cmd
TEST_DB=./matesays.db

.PHONY: all run fmt build clean test test-v test-race test-integration

all: fmt build

run:
	go run $(MAIN_PATH)

fmt:
	go fmt ./...

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(TEST_DB)

# Unit tests — runs all *_test.go files recursively
test:
	go test ./...

test-v:
	go test -v ./...

test-race:
	go test -race ./...

# Integration tests: exercise the compiled binary end-to-end
test-integration: build
	@echo "=== Integration tests ==="
	@rm -f $(TEST_DB)

	# Test: add first quote
	@echo "--- add first ---"
	$(BUILD_DIR)/$(BINARY_NAME) --add "something mate said"

	# Test: add second quote
	@echo "--- add second ---"
	$(BUILD_DIR)/$(BINARY_NAME) --add "another quote"

	# Test: list all quotes
	@echo "--- list ---"
	$(BUILD_DIR)/$(BINARY_NAME) --list

	# Test: get quote by ID
	@echo "--- get ---"
	$(BUILD_DIR)/$(BINARY_NAME) --get 1

	# Test: random quote
	@echo "--- random ---"
	$(BUILD_DIR)/$(BINARY_NAME) --random

	# Test: delete one quote
	@echo "--- delete ---"
	$(BUILD_DIR)/$(BINARY_NAME) --delete 1

	# Test: verify deletion
	@echo "--- list after delete ---"
	$(BUILD_DIR)/$(BINARY_NAME) --list

	# Test: delete all remaining quotes
	@echo "--- delete-all ---"
	$(BUILD_DIR)/$(BINARY_NAME) --delete-all

	# Test: verify empty
	@echo "--- list after delete-all ---"
	$(BUILD_DIR)/$(BINARY_NAME) --list || true

	# Clean up test database
	@rm -f $(TEST_DB)
	@echo "=== Integration tests complete ==="
