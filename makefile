# Define the name of the executable
BINARY_NAME=out

# Define the Go source file
SRC=.

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