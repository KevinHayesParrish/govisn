.PHONY: build clean test release release-build release-build-linux

# Read version from VERSION.txt file (skip comments)
VERSION := $(shell grep -v '^#' VERSION.txt | head -1)

# Build the application
build:
	go build -ldflags "-X main.GOVISN_VERSION=$(VERSION)" -o govisn

# Clean build artifacts
clean:
	rm -f govisn

# Run tests (if available)
test:
	go test ./...

# Release: tag and prepare for distribution
release: test
	@echo "Preparing release v$(VERSION)"
	git add VERSION.txt CHANGELOG.md
	git commit -m "Bump version to $(VERSION)" || true
	git tag -a v$(VERSION) -m "Release version $(VERSION)" || true
	@echo "Release v$(VERSION) prepared. Push with: git push origin v$(VERSION)"

# Build release binaries for multiple platforms
release-build: clean
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.GOVISN_VERSION=$(VERSION)" -o govisn-$(VERSION)-linux-amd64
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.GOVISN_VERSION=$(VERSION)" -o govisn-$(VERSION)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.GOVISN_VERSION=$(VERSION)" -o govisn-$(VERSION)-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.GOVISN_VERSION=$(VERSION)" -o govisn-$(VERSION)-windows-amd64.exe
	@echo "Release binaries created for v$(VERSION)"

# Build release binary for Linux only (avoids cross-compilation issues)
release-build-linux: clean
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.GOVISN_VERSION=$(VERSION)" -o govisn-$(VERSION)-linux-amd64
	@echo "Linux release binary created for v$(VERSION)"

# Help target
help:
	@echo "Available targets:"
	@echo "  make build               - Build the application"
	@echo "  make clean               - Remove build artifacts"
	@echo "  make test                - Run tests"
	@echo "  make release             - Prepare a new release (tag and commit)"
	@echo "  make release-build       - Build release binaries for multiple platforms"
	@echo "  make release-build-linux - Build release binary for Linux only"
	@echo "  make help                - Show this help message"
