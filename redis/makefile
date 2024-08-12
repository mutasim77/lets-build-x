.PHONY: build run clean test clear-db run-redis

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# Binary name
BINARY_NAME=redis

# Main package path
MAIN_PACKAGE=.

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)

# Run the project
run:
	$(GORUN) $(MAIN_PACKAGE)

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./tests

# Format all Go files
fmt:
	$(GOCMD) fmt ./...

# Remove the database file
clear-db:
	rm -f database.aof

# Build and run the Redis clone
run-redis: build
	./$(BINARY_NAME)
