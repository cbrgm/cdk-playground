GOOS            := linux
GOARCH          := amd64
GO              := GOFLAGS=-mod=vendor go
BUILD_DIR       := bin
PACKAGES        := $(shell $(GO) list ./... | grep -v node_modules)
BINARIES        := $(shell find . -name 'main.go' | grep -v node_modules | grep -v vendor | awk -F/ '{print "bin/" $$3}')

$(BUILD_DIR)/%:
	env GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -ldflags="-s -w" -o $@/main ./lambda/$*

.PHONY: test
test:
	$(GO) test -cover ./...

.PHONY: build
build: $(BINARIES)

.PHONY: vendor
vendor:
	rm -rf vendor/
	go mod vendor

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)/