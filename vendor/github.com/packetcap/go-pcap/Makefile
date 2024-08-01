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

BINDIR := ./dist
BIN := pcap
GOBINDIR ?= $(shell go env GOPATH)/bin
LOCALBIN := $(BINDIR)/$(BIN)-$(OS)-$(ARCH)
INSTALLBIN := $(GOBINDIR)/$(BIN)

.PHONY: build clean fmt test fmt-check lint golangci-lint

export GO111MODULE=on

LINTER ?= $(GOBINDIR)/golangci-lint
LINTER_VERSION ?= v1.23.3
GOFILES := $(shell find . -name '*.go' | grep -v go/pkg/mod)

$(BINDIR):
	mkdir -p $@

build: $(LOCALBIN) $(BIN)
$(LOCALBIN): $(BINDIR)
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o $@ ./cmd
$(BIN):
	@if [ "$(OS)" = "$(BUILDOS)" -a "$(ARCH)" = "$(BUILDARCH)" -a ! -e "$@" ]; then ln -s $(LOCALBIN) $@; fi

install: $(INSTALLBIN)
$(INSTALLBIN):
	CGO_ENABLED=0 go build -o $@

clean:
	@rm -rf $(BIN) $(BINDIR)

fmt:
	gofmt -w -s $(GOFILES)

fmt-check:
	@FMTOUT=$$(gofmt -l $(GOFILES)); \
	if [ -n "$${FMTOUT}" ]; then echo $${FMTOUT}; exit 1; fi

vet:
	go vet ./...

golangci-lint: $(LINTER)
$(LINTER):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBINDIR) $(LINTER_VERSION)

## Lint the files
lint: golangci-lint
	@$(LINTER) run ./...

test:
	go test ./...
