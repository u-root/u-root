#!/usr/bin/env bash

# because things are never simple.
# See https://github.com/codecov/example-go#caveat-multiple-files

set -e
echo "" > coverage.txt

# show the network configuration. This can help troubleshooting integration
# tests.
ip a

GO_TEST_OPTS=()
if [[ "$TRAVIS_GO_VERSION" =~ ^1.(9|10|11|12)$ ]]
then
    # We use fmt.Errorf with verb "%w" which appeared only in Go1.13.
    # So the code compiles and works with Go1.12, but error descriptions
    # looks uglier and it does not pass "vet" tests on Go<1.13.
    GO_TEST_OPTS+='-vet=off'
fi

for d in $(go list ./... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic ${GO_TEST_OPTS[@]} $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
    # integration tests
    go test -c -cover -tags=integration -race -covermode=atomic ${GO_TEST_OPTS[@]} $d
    testbin="./$(basename $d).test"
    # only run it if it was built - i.e. if there are integ tests
    test -x "${testbin}" && sudo "./${testbin}" -test.coverprofile=profile.out
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm -f profile.out
    fi
done
