# Define the name of the executable
BINARY_NAME=out

# Define the Go source file
SRC=.

# Define the Goose command and the migration directory
GOOSE_CMD=goose
MIGRATION_DIR=./sql/schema

ifneq ("$(wildcard .env)","")
    include .env
    export
endif

# Default target
.PHONY: all
all: build

# Build the Go application
.PHONY: build
build:
	@echo "Building the application..."
	go build -o $(BINARY_NAME) $(SRC)
	./out

# Clean up the build artifacts and run goose down migration
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@if [ -f $(BINARY_NAME) ]; then \
		rm -f $(BINARY_NAME); \
		echo "Removed executable: $(BINARY_NAME)"; \
	fi

# Help target to display available commands
.PHONY: help
help:
	@echo "Makefile Commands:"
	@echo "  make build  - Build the application"
	@echo "  make run    - Run the application"
	@echo "  make clean  - Clean up"
	@echo "  make help   - Show this help message"

# Database UP migrations
.PHONE: db
db:
	@echo "Running goose up migration..."
	$(GOOSE_CMD) postgres $(DB_URL) up -dir $(MIGRATION_DIR)
	@echo "Generating go code from sql..."
	sqlc generate

# Database UP migrations
.PHONE: db_clean
db_clean:
	@echo "Running goose down migration..."
	$(GOOSE_CMD) postgres $(DB_URL) down -dir $(MIGRATION_DIR)