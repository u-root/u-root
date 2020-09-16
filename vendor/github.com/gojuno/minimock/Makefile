all: test lint

generate:
	go run ./cmd/minimock/minimock.go -i github.com/gojuno/minimock.Tester -o ./tests -s _mock_test.go
	go run ./cmd/minimock/minimock.go -i tests.Formatter -o ./tests -s _mock.go

lint:
	golint ./... && go vet ./...

test: generate
	go test -race ./...
