SHELL := /bin/bash 

# 'make' without arguments will show help
.DEFAULTGOAL := help

BINARY:=adm
BIN_DIR:=./bin
ROOT_DIR:=$(shell dirname $(MAKEFILE_LIST) | xargs)
TEST_PATH := $(shell sed -e 's/ /\\\ /g' <<< "$(ROOT_DIR)/tests")
TEST_PACKAGES := "args,sources,libadm/loaders,libadm/model,libadm/graph,libadm/graphviz,libadm/export"

help: # Show this help
	@egrep -h '\s#\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

all: test build clean

build: # Build program
	@echo -n "Checking if './bin' directory exists..."
	@if [ ! -d $(BIN_DIR) ]; then mkdir $(BIN_DIR); echo "CREATED"; else echo "YES"; fi

	@echo -n "Building 'bin/adm'..."
	@cd src/main && go get args
	@cd src/main && go build -ldflags "-s -w" -v -o ../../$(BIN_DIR)/$(BINARY) .
	@echo "DONE"

test: # Run automated tests
	@echo -n "Running automated tests..."

	@cd $(TEST_PATH) && go get args && go get -t tests
	@cd $(TEST_PATH) && go test -v -cover -coverpkg $(TEST_PACKAGES) --coverprofile /tmp/adm-cli.coverage.out ./...
	@sed -E "s/(^.*\/.*:)/.\/src\/\1/g" /tmp/adm-cli.coverage.out > /tmp/adm-cli.coverage.fullpath.out

report: # Coverage report (after `make tests`)
	@go tool cover -html=/tmp/adm-cli.coverage.fullpath.out

	@cd $(TEST_PATH) && go get args && go get -t tests
	@cd $(TEST_PATH) && go test -cpuprofile /tmp/adm-cli.cpu.prof -memprofile /tmp/adm-cli.mem.prof -bench .
	@echo "=================================="
	@echo "CPU Profiler results..."
	@go tool pprof -text /tmp/adm-cli.cpu.prof
	@echo "=================================="
	@echo "Memory Profiler results..."
	@go tool pprof -text /tmp/adm-cli.mem.prof
	@echo "=================================="

clean: # Clean all build/test artifacts.
	@echo "Removing build/test artifacts..."
	@rm -f $(BIN_DIR)/$(BINARY)
	@rm -f $(BIN_DIR)/graphviz.svg
	@rm -f $(BIN_DIR)/export/*.attack
	@rm -f $(BIN_DIR)/export/*.feature
	@rm -rf $(BIN_DIR)/export
	@rm -rf $(BIN_DIR)
	@echo "DONE"