# BUILDARCH is the host architecture
# ARCH is the target architecture
# we need to keep track of them separately
BUILDARCH ?= $(shell uname -m)
BUILDOS ?= $(shell uname -s | tr A-Z a-z)

# canonicalized names for host architecture
ifeq ($(BUILDARCH),aarch64)
BUILDARCH=arm64
endif
ifeq ($(BUILDARCH),x86_64)
BUILDARCH=amd64
endif

# unless otherwise set, I am building for my own architecture, i.e. not cross-compiling
# and for my OS
ARCH ?= $(BUILDARCH)
OS ?= $(BUILDOS)

# canonicalized names for target architecture
ifeq ($(ARCH),aarch64)
        override ARCH=arm64
endif
ifeq ($(ARCH),x86_64)
    override ARCH=amd64
endif

# The default build creates the pcap binary in the project root.
# Set BINDIR to place build artifacts elsewhere when needed.
BINDIR ?= .
BIN ?= pcap
# TARGET labels an artifact independently of the Go target tuple. Keep the
# supported 32-bit ARM baselines explicit so a v7 binary is never mislabeled
# as suitable for v6.
TARGETARCH := $(ARCH)
ifeq ($(ARCH),arm)
ifeq ($(filter 6 7,$(GOARM)),)
$(error GOARM must be 6 or 7 when ARCH=arm)
endif
TARGETARCH := armv$(GOARM)
endif
TARGET ?= $(OS)-$(TARGETARCH)
GOBINDIR ?= $(shell go env GOPATH)/bin
HOSTBIN := $(BINDIR)/$(BIN)
TARGETBIN := $(BINDIR)/$(BIN)-$(TARGET)
ifeq ($(OS),$(BUILDOS))
ifeq ($(ARCH),$(BUILDARCH))
LOCALBIN := $(HOSTBIN)
else
LOCALBIN := $(TARGETBIN)
endif
else
LOCALBIN := $(TARGETBIN)
endif
INSTALLBIN := $(GOBINDIR)/$(BIN)

.PHONY: build build-artifact clean fmt test integration bench fmt-check lint golangci-lint format-tools

export GO111MODULE=on

LINTER ?= $(GOBINDIR)/golangci-lint
LINTER_VERSION ?= v2.5.0
GOIMPORTS ?= $(GOBINDIR)/goimports
GOIMPORTS_VERSION ?= v0.38.0
GOFUMPT ?= $(GOBINDIR)/gofumpt
GOFUMPT_VERSION ?= v0.9.1
SHFMT ?= $(GOBINDIR)/shfmt
SHFMT_VERSION ?= v3.12.0
MODULE_PATH := github.com/huatuo-ai/go-pcap
GOFILES := $(shell git ls-files '*.go')
SHELLFILES := $(shell git ls-files '*.sh')

$(BINDIR):
	mkdir -p $@

build: $(LOCALBIN)

# build-artifact always includes the target suffix. Release builds use this
# target so a native build cannot overwrite the target-specific artifact name.
build-artifact: $(TARGETBIN)

$(HOSTBIN) $(TARGETBIN): $(BINDIR)
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) $(if $(filter arm,$(ARCH)),$(if $(GOARM),GOARM=$(GOARM))) go build -o $@ ./cmd

install: $(INSTALLBIN)
$(INSTALLBIN):
	CGO_ENABLED=0 go build -o $@

clean:
	@rm -f $(HOSTBIN) $(TARGETBIN)

format-tools:
	@GOBIN=$(GOBINDIR) go install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)
	@GOBIN=$(GOBINDIR) go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
	@GOBIN=$(GOBINDIR) go install mvdan.cc/sh/v3/cmd/shfmt@$(SHFMT_VERSION)

fmt: format-tools
	$(GOIMPORTS) -w -local $(MODULE_PATH) $(GOFILES)
	$(GOFUMPT) -w $(GOFILES)
	gofmt -w -r 'interface{} -> any' $(GOFILES)
	$(SHFMT) -i 0 -w $(SHELLFILES)

fmt-check: format-tools
	@FMTOUT=$$($(GOIMPORTS) -l -local $(MODULE_PATH) $(GOFILES)); \
	if [ -n "$${FMTOUT}" ]; then echo "$${FMTOUT}"; exit 1; fi
	@FMTOUT=$$($(GOFUMPT) -l $(GOFILES)); \
	if [ -n "$${FMTOUT}" ]; then echo "$${FMTOUT}"; exit 1; fi
	@FMTOUT=$$(gofmt -l -r 'interface{} -> any' $(GOFILES)); \
	if [ -n "$${FMTOUT}" ]; then echo "$${FMTOUT}"; exit 1; fi
	@$(SHFMT) -i 0 -d $(SHELLFILES)

vet:
	go vet ./...

golangci-lint:
	@if ! $(LINTER) version 2>/dev/null | grep -q "version $(LINTER_VERSION)"; then \
		GOBIN=$(GOBINDIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(LINTER_VERSION); \
	fi

## Lint the files
lint: golangci-lint
	@$(LINTER) run ./...

test:
	go test ./...

integration: build
	@bash test/run.sh

bench:
	go test ./filter -run '^$$' -bench . -benchmem -count=10
