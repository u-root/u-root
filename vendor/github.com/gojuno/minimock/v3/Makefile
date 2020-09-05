all: install test lint clean

generate:
	go run ./cmd/minimock/minimock.go -i github.com/gojuno/minimock/v3.Tester -o ./tests
	go run ./cmd/minimock/minimock.go -i ./tests.Formatter -o ./tests/formatter_mock.go

lint:
	gometalinter ./tests/ -I minimock -e gopathwalk --disable=gotype --deadline=2m

install:
	go mod download
	go install ./cmd/minimock

clean:
	[ -e ./tests/formatter_mock.go.test_origin ] && mv -f ./tests/formatter_mock.go.test_origin ./tests/formatter_mock.go
	[ -e ./tests/tester_mock_test.go.test_origin ] && mv -f ./tests/tester_mock_test.go.test_origin ./tests/tester_mock_test.go
	rm -Rf bin/ dist/

test_save_origin:
	[ -e ./tests/formatter_mock.go.test_origin ] || cp ./tests/formatter_mock.go ./tests/formatter_mock.go.test_origin
	[ -e ./tests/tester_mock_test.go.test_origin ] || cp ./tests/tester_mock_test.go ./tests/tester_mock_test.go.test_origin

test: test_save_origin generate
	diff ./tests/formatter_mock.go ./tests/formatter_mock.go.test_origin
	diff ./tests/tester_mock_test.go ./tests/tester_mock_test.go.test_origin
	go test -race ./...
