.DEFAULT_GOAL := help
 
BIN_NAME := web3datacli
BUILD_DIR := ./out

GIT_SHA := $(shell git rev-parse HEAD | cut -c 1-8)
GIT_TAG := $(shell git describe --tags)
DATE := $(shell date +%s)
VERSION_FLAGS=\
  -X github.com/thewhitewizard/web3data-cli/cmd/version.Version=$(GIT_TAG) \
  -X github.com/thewhitewizard/web3data-cli/cmd/version.Commit=$(GIT_SHA) \
  -X github.com/thewhitewizard/web3data-cli/cmd/version.Date=$(DATE) \
  -X github.com/thewhitewizard/web3data-cli/cmd/version.BuiltBy=makefile

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: $(BUILD_DIR)
$(BUILD_DIR): ## Create the build folder.
	mkdir -p $(BUILD_DIR)

.PHONY: build
build: $(BUILD_DIR) ## Build go binary.
	go build -ldflags "$(VERSION_FLAGS)" -o $(BUILD_DIR)/$(BIN_NAME) main.go

.PHONY: cross
cross: $(BUILD_DIR) ## Cross-compile go binaries without using CGO.	
	GOOS=linux  GOARCH=amd64 go build -o $(BUILD_DIR)/$(BIN_NAME)_$(GIT_TAG)_linux_amd64  main.go
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BIN_NAME)_$(GIT_TAG)_darwin_amd64 main.go

.PHONY: clean
clean: ## Clean the binary folder.
	$(RM) -r $(BUILD_DIR)
