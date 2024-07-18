# Define Go command
GO_CMD := go

# Define build directory
BUILD_DIR := ./build

# Define tools directory
TOOLS_DIR := ./tools

# Program definitions
ALGODHEALTH_SRC := $(TOOLS_DIR)/algodhealth/algodhealth.go
CATCH_CATCHPOINT_SRC := $(TOOLS_DIR)/catch-catchpoint/catch-catchpoint.go
GET_METRICS_SRC := $(TOOLS_DIR)/get-metrics/get-metrics.go
START_NODE_SRC := $(TOOLS_DIR)/start-node/start-node.go
START_METRICS_SRC := $(TOOLS_DIR)/start-metrics/start-metrics.go

ALGODHEALTH_OUT := $(BUILD_DIR)/algodhealth
CATCH_CATCHPOINT_OUT := $(BUILD_DIR)/catch-catchpoint
GET_METRICS_OUT := $(BUILD_DIR)/get-metrics
START_NODE_OUT := $(BUILD_DIR)/start-node
START_METRICS_OUT := $(BUILD_DIR)/start-metrics

TESTNET := testnet
NETWORK :=

# Default target
all: algodhealth catch-catchpoint get-metrics start-node start-metrics

testnet: NETWORK := $(TESTNET)
testnet: all

# Build targets
algodhealth:
	@mkdir -p $(BUILD_DIR)
	$(GO_CMD) build -o $(ALGODHEALTH_OUT) $(ALGODHEALTH_SRC)

catch-catchpoint:
	@mkdir -p $(BUILD_DIR)
	$(GO_CMD) build -o $(CATCH_CATCHPOINT_OUT) $(CATCH_CATCHPOINT_SRC)

get-metrics:
	@mkdir -p $(BUILD_DIR)
	$(GO_CMD) build -o $(GET_METRICS_OUT) $(GET_METRICS_SRC)

start-node:
	@mkdir -p $(BUILD_DIR)
	$(GO_CMD) build -o $(START_NODE_OUT) $(START_NODE_SRC)

start-metrics:
	@mkdir -p $(BUILD_DIR)
	$(GO_CMD) build -o $(START_METRICS_OUT) $(START_METRICS_SRC)

# Test target
test:
	@$(GO_CMD) test -v ./...

# Clean target
clean:
	@rm -rf $(BUILD_DIR)

.PHONY: all algodhealth catch-catchpoint get-metrics start-node start-metrics clean testnet