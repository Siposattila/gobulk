VERSION := $(shell echo $(VERSION))
OUTPUT := gobulk
BUILD_DIR := build
DIST_DIR := dist
TAR_FILE := $(DIST_DIR)/$(OUTPUT)_$(VERSION).tar.gz

all: build compress

build:
	@echo "Building GoBulk..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X main.buildTime=`date -u '+%Y.%m.%d.'` -X main.buildHash=`git rev-parse HEAD` -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(OUTPUT) .
	@echo "Copying .html files..."
	@cp *.html $(BUILD_DIR) 2>/dev/null || true

compress: build
	@echo "Compressing the build output..."
	@mkdir -p $(DIST_DIR)
	@tar -czvf $(TAR_FILE) -C $(BUILD_DIR) $(OUTPUT) *.html
	@echo "Compressed file created at $(TAR_FILE)"

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "Clean up done."

.PHONY: all build compress clean
