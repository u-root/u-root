export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)
export GOFLAGS := -mod=mod

all: install test lint

generate:
	go run ./cmd/minimock/minimock.go -i github.com/gojuno/minimock/v3.Tester -o ./tests
	go run ./cmd/minimock/minimock.go -i ./tests.Formatter -o ./tests/formatter_mock.go
	go run ./cmd/minimock/minimock.go -i ./tests.Formatter -o ./tests/formatter_with_custom_name_mock.go -n CustomFormatterNameMock
	go run ./cmd/minimock/minimock.go -i ./tests.genericInout -o ./tests/generic_inout.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericOut -o ./tests/generic_out.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericIn -o ./tests/generic_in.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericSpecific -o ./tests/generic_specific.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericSimpleUnion -o ./tests/generic_simple_union.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericComplexUnion -o ./tests/generic_complex_union.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericInlineUnion -o ./tests/generic_inline_union.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericInlineUnionWithManyTypes -o ./tests/generic_inline_with_many_options.go
	go run ./cmd/minimock/minimock.go -i ./tests.genericMultipleTypes -o ./tests/generic_multiple_args_with_different_types.go
	go run ./cmd/minimock/minimock.go -i ./tests.contextAccepter -o ./tests/context_accepter_mock.go
	go run ./cmd/minimock/minimock.go -i github.com/gojuno/minimock/v3.Tester -o ./tests/package_name_specified_test.go -p tests_test
	go run ./cmd/minimock/minimock.go -i ./tests.actor -o ./tests/actor_mock.go
	go run ./cmd/minimock/minimock.go -i ./tests.formatterAlias -o ./tests/formatter_alias_mock.go
	go run ./cmd/minimock/minimock.go -i ./tests.formatterType -o ./tests/formatter_type_mock.go
	go run ./cmd/minimock/minimock.go -i ./tests.reader -o ./tests/reader_mock.go -gr

./bin:
	mkdir ./bin

lint: ./bin/golangci-lint
	./bin/golangci-lint run --enable=goimports --disable=unused --exclude=S1023,"Error return value" ./tests/...

install:
	go mod download
	go install ./cmd/minimock

./bin/golangci-lint: ./bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.64.8

./bin/goreleaser: ./bin
	go install -modfile tools/go.mod github.com/goreleaser/goreleaser

./bin/minimock:
	go build ./cmd/minimock -o ./bin/minimock

.PHONY:
test:
	go test -race ./... -v

build: ./bin/goreleaser
	./bin/goreleaser build --snapshot --rm-dist

.PHONY:
tidy:
	go mod tidy
	cd tools && go mod tidy
