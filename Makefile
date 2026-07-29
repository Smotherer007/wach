APP_NAME    := Wach
APP_BUNDLE  := $(APP_NAME).app
BIN_DIR     := bin
BUILD_DIR   := $(BIN_DIR)/$(APP_BUNDLE)/Contents
RES_DIR     := $(BUILD_DIR)/Resources
MACOS_DIR   := $(BUILD_DIR)/MacOS
BINARY      := $(MACOS_DIR)/wach
GO_FLAGS    := -ldflags="-s -w"

.PHONY: all build install uninstall run clean test

all: build

build: clean
	@echo "==> Wach bauen..."
	@mkdir -p $(RES_DIR) $(MACOS_DIR)
	@cp appInfo/Info.plist $(BUILD_DIR)/
	@cp appInfo/icon.icns $(RES_DIR)/
	CGO_ENABLED=1 go build $(GO_FLAGS) -o $(BINARY) ./cmd/
	@echo "  ✅ $(APP_BUNDLE) gebaut in $(BIN_DIR)/"

install: build
	@echo "==> Installieren nach /Applications/..."
	@rm -rf /Applications/$(APP_BUNDLE)
	@cp -r $(BIN_DIR)/$(APP_BUNDLE) /Applications/
	@echo "  ✅ Installiert"

uninstall:
	-pkill -f "$(APP_BUNDLE)" 2>/dev/null || true
	@sleep 1
	@rm -rf /Applications/$(APP_BUNDLE)

run:
	CGO_ENABLED=1 go run ./cmd/

clean:
	@rm -rf $(BIN_DIR) 2>/dev/null || true

test:
	CGO_ENABLED=1 go test -v -race -coverprofile=cover.out ./...

coverage: test
	go tool cover -html=cover.out -o cover.html
