VENDOR := vendor
PKGS := $(shell go list ./... | grep -v /$(VENDOR)/)
SRC = $(shell find . -type f -name '*.go' -not -path "*/$(VENDOR)/*")
BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

.PHONY: test
test: vet
	go test $(PKGS)

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	$(GOMETALINTER) --install &> /dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	$(GOMETALINTER) ./... --vendor

.PHONY: fmt
fmt:
	@echo "Running gofmt on all sources..."
	@gofmt -s -l -w $(SRC)

.PHONY: fmtcheck
fmtcheck:
	@bash -c "diff -u <(echo -n) <(gofmt -d $(SRC))"

.PHONY: vet
vet:
	go vet $(PKGS)

.PHONY: cover
cover:
	$(shell [ -e coverage.out ] && rm coverage.out)
	@echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PKGS),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out -o=coverage-all.html
