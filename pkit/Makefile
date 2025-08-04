VERSION := 0.1.0
BINARY := pkit
BUILD_DIR := build

.PHONY: all clean build-linux build-darwin build-windows version-bump

all: clean build-linux build-darwin build-windows

clean:
	rm -rf $(BUILD_DIR)

build-linux:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY)-linux-amd64

build-darwin:
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY)-darwin-amd64

build-windows:
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe

version-bump:
	@echo "Current version: $(VERSION)"
	@read -p "Enter new version: " NEW_VERSION; \
	sed -i.bak "s/VERSION := .*/VERSION := $$NEW_VERSION/" Makefile && \
	sed -i.bak "s/version: .*/version: $$NEW_VERSION/" config.yaml && \
	rm -f Makefile.bak config.yaml.bak && \
	echo "Version updated to $$NEW_VERSION"

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY)