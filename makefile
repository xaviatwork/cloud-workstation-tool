# --- Configuration ---

# The name of the resulting binary/application
APP_NAME := cw-tunnel
# Directory to place the compiled binary
BUILD_DIR := release
# The main package path relative to the root (e.g., ".")
MAIN_PACKAGE := cmd/tunnel/*.go
# Target OS and Architecture for Apple Silicon
TARGET_MAC_OS := darwin
TARGET_MAC_ARCH := arm64
# Target OS and Architecture for Intel Windows
TARGET_OS_WIN := windows
TARGET_ARCH_WIN := amd64
# The final name of the output file
OUTPUT_NAME := $(APP_NAME)-$(TARGET_OS_WIN)-$(TARGET_ARCH_WIN)

# --- Targets ---
.PHONY: all build clean

# Default target: builds the macOS ARM binary
all: mac

# 1. Build Target
mac:
	@echo "--- Compiling $(APP_NAME) for $(TARGET_MAC_OS)/$(TARGET_MAC_ARCH) ---"
	@# Execute the cross-compilation command
	GOOS=$(TARGET_MAC_OS) GOARCH=$(TARGET_MAC_ARCH) go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PACKAGE)
	@echo "--> Built: $(BUILD_DIR)/$(APP_NAME)"

windows:
	@echo "--- Compiling $(APP_NAME) for $(TARGET_OS_WIN)/$(TARGET_ARCH_WIN) ---"
	@# Execute the cross-compilation command
	GOOS=$(TARGET_OS_WIN) GOARCH=$(TARGET_ARCH_WIN) go build -o $(BUILD_DIR)/$(OUTPUT_NAME) $(MAIN_PACKAGE)
	@echo "   -> Built: $(BUILD_DIR)/$(OUTPUT_NAME)"
