APP_NAME    := Wach
APP_BUNDLE  := $(APP_NAME).app
BIN_DIR     := bin
BUILD_DIR   := $(BIN_DIR)/$(APP_BUNDLE)/Contents
RES_DIR     := $(BUILD_DIR)/Resources
MACOS_DIR   := $(BUILD_DIR)/MacOS
BINARY      := $(MACOS_DIR)/wach
GO_FLAGS    := -ldflags="-s -w"
GO_FLAGS_ARM := CGO_ENABLED=1 GOARCH=arm64 GOOS=darwin
GO_FLAGS_INTEL := CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin

.PHONY: all build build-intel install uninstall run clean test lint

all: build

# --- Apple Silicon Build (default) ---
build: clean
	@echo "==> Wach bauen (Apple Silicon)..."
	@mkdir -p $(RES_DIR) $(MACOS_DIR)
	@cp appInfo/Info.plist $(BUILD_DIR)/
	@cp appInfo/icon.icns $(RES_DIR)/
	@$(GO_FLAGS_ARM) go build $(GO_FLAGS) -o $(BINARY) ./cmd/
	@echo ""
	@echo "  ✅ $(APP_BUNDLE) gebaut!"
	@echo "  📁 $(BIN_DIR)/$(APP_BUNDLE)"
	@echo "  🔧 natives arm64 Binary"
	@echo ""
	@echo "  Installieren:  make install"
	@echo "  Ausfuehren:    make run"

# --- Intel Build ---
build-intel: clean
	@echo "==> Wach bauen (Intel)..."
	@mkdir -p $(RES_DIR) $(MACOS_DIR)
	@cp appInfo/Info.plist $(BUILD_DIR)/
	@cp appInfo/icon.icns $(RES_DIR)/
	@$(GO_FLAGS_INTEL) go build $(GO_FLAGS) -o $(BINARY) ./cmd/
	@echo ""
	@echo "  ✅ $(APP_BUNDLE) (Intel) gebaut!"

# --- Install ---
install: build
	@echo "==> Installiere nach /Applications/..."
	@rm -rf /Applications/$(APP_BUNDLE)
	@cp -r $(BIN_DIR)/$(APP_BUNDLE) /Applications/
	@echo "  ✅ Installiert: /Applications/$(APP_BUNDLE)"
	@echo ""
	@echo "  Wichtig: Beim ersten Start Rechtsklick > Offnen wahlen."
	@echo "  Danach Zuganglichkeit erlauben unter:"
	@echo "  Systemeinstellungen > Datenschutz > Bedienungshilfen > ✓ Wach"

# --- Uninstall ---
uninstall:
	@echo "==> Deinstalliere..."
	-pkill -f "$(APP_BUNDLE)" 2>/dev/null || true
	@sleep 1
	@rm -rf /Applications/$(APP_BUNDLE)
	@echo "  ✅ Deinstalliert"

# --- Run directly (no bundle) ---
run:
	@$(GO_FLAGS_ARM) go run ./cmd/

# --- Clean ---
clean:
	@rm -rf $(BIN_DIR) 2>/dev/null || true

# --- Tests ---
test:
	@CGO_ENABLED=1 go test -v -race -coverprofile=cover.out ./...

coverage: test
	@go tool cover -html=cover.out -o cover.html
	@echo "Coverage report: cover.html"

# --- Lint ---
lint:
	@which staticcheck 2>/dev/null || go install honnef.co/go/tools/cmd/staticcheck@latest
	@staticcheck ./...
