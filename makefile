# --- Configuration ---

# The name of the resulting binary/application
APP_NAME := cw
# Directory to place the compiled binary
BUILD_DIR := release
# The main package path relative to the root (e.g., ".")
MAIN_PACKAGE := *.go
# Target OS
TARGET_OS_MAC := darwin
TARGET_OS_WIN := windows
TARGET_OS_LINUX := linux
# Target architectures
TARGET_ARCH_ARM64 := arm64
TARGET_ARCH_AMD64 := amd64

# --- Targets ---
.PHONY: all mac windows linux

all: mac windows linux

mac:
	@echo "--- Compiling $(APP_NAME) for $(TARGET_OS_MAC)/$(TARGET_ARCH_ARM64) ---"
	@# Execute the cross-compilation command
	GOOS=$(TARGET_OS_MAC) GOARCH=$(TARGET_ARCH_ARM64) go build -o $(BUILD_DIR)/$(APP_NAME)-$(TARGET_OS_MAC)-$(TARGET_ARCH_ARM64) $(MAIN_PACKAGE)
	@echo "--> Built: $(BUILD_DIR)/$(APP_NAME)-$(TARGET_OS_MAC)-$(TARGET_ARCH_ARM64)"

	@echo "--- Compiling $(APP_NAME) for $(TARGET_OS_MAC)/$(TARGET_ARCH_AMD64) ---"
	@# Execute the cross-compilation command
	GOOS=$(TARGET_MAC_OS) GOARCH=$(TARGET_ARCH_AMD64) go build -o $(BUILD_DIR)/$(APP_NAME)-$(TARGET_OS_MAC)-$(TARGET_ARCH_AMD64) $(MAIN_PACKAGE)
	@echo "--> Built: $(BUILD_DIR)/$(APP_NAME)-$(TARGET_OS_MAC)-$(TARGET_ARCH_AMD64)"

windows:
	@echo "--- Compiling $(APP_NAME) for $(TARGET_OS_WIN)/$(TARGET_ARCH_AMD64) ---"
	@# Execute the cross-compilation command
	GOOS=$(TARGET_OS_WIN) GOARCH=$(TARGET_ARCH_AMD64) go build -o $(BUILD_DIR)/$(APP_NAME)-$(TARGET_OS_WIN)-$(TARGET_ARCH_AMD64) $(MAIN_PACKAGE)
	@echo "   -> Built: $(BUILD_DIR)/$(APP_NAME)-$(TARGET_OS_WIN)-$(TARGET_ARCH_AMD64)"

linux:
	@echo "--- Compiling $(APP_NAME) for $(TARGET_OS_LINUX)/$(TARGET_ARCH_AMD64) ---"
	@# Execute the cross-compilation command
	GOOS=$(TARGET_OS_LINUX) GOARCH=$(TARGET_ARCH_AMD64) go build -o $(BUILD_DIR)/$(APP_NAME)-$(TARGET_OS_LINUX)-$(TARGET_ARCH_AMD64) $(MAIN_PACKAGE)
	@echo "   -> Built: $(BUILD_DIR)/$(APP_NAME)-$(TARGET_OS_LINUX)-$(TARGET_ARCH_AMD64)"