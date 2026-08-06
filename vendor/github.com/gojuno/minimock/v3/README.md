![logo](https://rawgit.com/gojuno/minimock/master/logo.svg)
[![GoDoc](https://godoc.org/github.com/gojuno/minimock?status.svg)](http://godoc.org/github.com/gojuno/minimock) 
[![Go Report Card](https://goreportcard.com/badge/github.com/gojuno/minimock)](https://goreportcard.com/report/github.com/gojuno/minimock)
[![Release](https://img.shields.io/github/release/gojuno/minimock.svg)](https://github.com/gojuno/minimock/releases/latest)
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#testing)


## Summary 
Minimock generates mocks out of Go interface declarations.

The main features of minimock are:

* It generates statically typed mocks and helpers. There's no need for type assertions when you use minimock.
* It's fully integrated with the standard Go "testing" package.
* It's ready for Go modules.
* It supports generics.
* It works well with [table driven tests](https://dave.cheney.net/2013/06/09/writing-table-driven-tests-in-go) because you can set up mocks for several methods in one line of code using the builder pattern.
* It can generate several mocks in one run.
* It can generate mocks from interface aliases.
* It generates code that passes default set of [golangci-lint](https://github.com/golangci/golangci-lint) checks.
* It puts //go:generate instruction into the generated code, so all you need to do when the source interface is updated is to run the `go generate ./...` command from within the project's directory.
* It makes sure that all mocked methods have been called during the test and keeps your test code clean and up to date.
* It provides When and Then helpers to set up several expectations and results for any method.
* It generates concurrent-safe mocks and mock invocation counters that you can use to manage mock behavior depending on the number of calls.
* It can be used with the [GoUnit](https://github.com/hexdigest/gounit) tool which generates table-driven tests that make use of minimock.

## Installation

If you use go modules please download the [latest binary](https://github.com/gojuno/minimock/releases/latest)
or install minimock from source:
```
go install github.com/gojuno/minimock/v3/cmd/minimock@latest
```

If you don't use go modules please find the latest v2.x binary [here](https://github.com/gojuno/minimock/releases)
or install minimock using [v2 branch](https://github.com/gojuno/minimock/tree/v2)

## Usage

```
 minimock [-i source.interface] [-o output/dir/or/file.go] [-g]
  -g	don't put go:generate instruction into the generated code
  -h	show this help message
  -i string
    	comma-separated names of the interfaces to mock, i.e fmt.Stringer,io.Reader
    	use io.* notation to generate mocks for all interfaces in the "io" package (default "*")
  -o string
    	comma-separated destination file names or packages to put the generated mocks in,
    	by default the generated mock is placed in the source package directory
  -p string 
        comma-separated package names,
        by default the generated package names are taken from the destination directory names
  -pr string
        mock file prefix
  -s string
    	mock file suffix (default "_mock_test.go")
  -gr
        changes go:generate line from "//go:generate minimock args..." to  
        "//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock", 
        useful while controlling minimock version with go mod
```

Let's say we have the following interface declaration in github.com/gojuno/minimock/tests package:
```go
type Formatter interface {
	Format(string, ...interface{}) string
}
```

This will generate mocks for all interfaces defined in the "tests" package:

```
$ cd ~/go/src/github.com/gojuno/minimock/tests
$ minimock 
```

Here is how to generate a mock for the "Formatter" interface only:

```
$ cd ~/go/src/github.com/gojuno/minimock/tests
$ minimock -i Formatter 
```

Same using the relative package notation:

```
$ minimock -i ./tests.Formatter
```

Same using the full import path of the source package:

```
$ minimock -i github.com/gojuno/minimock/tests.Formatter -o ./tests/
```

All the examples above generate ./tests/formatter_mock_test.go file


Now it's time to use the generated mock. There are several ways it can be done.

### Setting up a mock using the builder pattern and Expect/Return methods:
```go
mc := minimock.NewController(t)
formatterMock := NewFormatterMock(mc).FormatMock.Expect("hello %s!", "world").Return("hello world!")
```

The builder pattern is convenient when you have more than one method to mock.
Let's say we have an io.ReadCloser interface which has two methods: Read and Close
```go
type ReadCloser interface {
	Read(p []byte) (n int, err error)
	Close() error
}
```

We can set up a mock using a simple one-liner:
```go
mc := minimock.NewController(t)
readCloserMock := NewReadCloserMock(mc).ReadMock.Expect([]byte(1,2,3)).Return(3, nil).CloseMock.Return(nil)
```

But what if we don't want to check all arguments of the read method?
Let's say we just want to check that the second element of the given slice "p" is 2.
This is where "Inspect" helper comes into play:
```go
mc := minimock.NewController(t)
readCloserMock := NewReadCloserMock(mc).ReadMock.Inspect(func(p []byte){
  assert.Equal(mc, 2, p[1])
}).Return(3, nil).CloseMock.Return(nil)

```

### Setting up a mock using ExpectParams helpers:

Let's say we have a mocking interface with function that has many arguments:
```go
type If interface {
	Do(intArg int, stringArg string, floatArg float)
}
```

Imagine that you don't want to check all the arguments, just one or two of them.
Then we can use ExpectParams helpers, which are generated for each argument:
```go
mc := minimock.NewController(t)
ifMock := NewIfMock(mc).DoMock.ExpectIntArgParam1(10).ExpectFloatArgParam3(10.2).Return()
```

### Setting up a mock using When/Then helpers:
```go
mc := minimock.NewController(t)
formatterMock := NewFormatterMock(mc)
formatterMock.FormatMock.When("Hello %s!", "world").Then("Hello world!")
formatterMock.FormatMock.When("Hi %s!", "there").Then("Hi there!")
```

alternatively you can use the one-liner:

```go
formatterMock = NewFormatterMock(mc).FormatMock.When("Hello %s!", "world").Then("Hello world!").FormatMock.When("Hi %s!", "there").Then("Hi there!")
```

### Setting up a mock using the Set method:
```go
mc := minimock.NewController(t)
formatterMock := NewFormatterMock(mc).FormatMock.Set(func(string, ...interface{}) string {
  return "minimock"
})
```

You can also use invocation counters in your mocks and tests:
```go
mc := minimock.NewController(t)
formatterMock := NewFormatterMock(mc)
formatterMock.FormatMock.Set(func(string, ...interface{}) string {
  return fmt.Sprintf("minimock: %d", formatterMock.BeforeFormatCounter())
})
```

### Setting up expected times mock was called:
Imagine you expect mock to be called exactly 10 times. 
Then you can set `Times` helper to check how many times mock was invoked.

```go
mc := minimock.NewController(t)
formatterMock := NewFormatterMock(mc).FormatMock.Times(10).Expect("hello %s!", "world").Return("hello world!")
```

There are also cases, when you don't know for sure if the mocking method would be called or not. 
But you still want to mock it, if it will be called. This is where "Optional" option comes into play:

```go
mc := minimock.NewController(t)
formatterMock := NewFormatterMock(mc).FormatMock.Optional().Expect("hello %s!", "world").Return("hello world!")
```

When this option is set, it disables checking the call of mocking method. 

### Mocking context
Sometimes context gets modified by the time the mocked method is being called.
However, in most cases you don't really care about the exact value of the context argument.
In such cases you can use special `minimock.AnyContext` variable, here are a couple of examples:

```go
mc := minimock.NewController(t)
senderMock := NewSenderMock(mc).
  SendMock.
    When(minimock.AnyContext, "message1").Then(nil).
    When(minimock.AnyContext, "message2").Then(errors.New("invalid message"))
```

or using Expect:

```go
mc := minimock.NewController(t)
senderMock := NewSenderMock(mc).
  SendMock.Expect(minimock.AnyContext, "message").Return(nil)
```

### Make sure that your mocks are being used 
Often we write tons of mocks to test our code but sometimes the tested code stops using mocked dependencies.
You can easily identify this problem by using `minimock.NewController` instead of just `*testing.T`. 
Alternatively you can use `mc.Wait` helper if your're testing concurrent code.
These helpers ensure that all your mocks and expectations have been used at least once during the test run.

```go
func TestSomething(t *testing.T) {
  // it will mark this example test as failed because there are no calls
  // to formatterMock.Format() and readCloserMock.Read() below
  mc := minimock.NewController(t)

  formatterMock := NewFormatterMock(mc)
  formatterMock.FormatMock.Return("minimock")

  readCloserMock := NewReadCloserMock(mc)
  readCloserMock.ReadMock.Return(5, nil)
}
```

### Testing concurrent code
Testing concurrent code is tough. Fortunately minimock.Controller provides you with the helper method that makes testing concurrent code easy.
Here is how it works:

```go
func TestSomething(t *testing.T) {
  mc := minimock.NewController(t)

  //Wait ensures that all mocked methods have been called within the given time span
  //if any of the mocked methods have not been called Wait marks the test as failed
  defer mc.Wait(time.Second)

  formatterMock := NewFormatterMock(mc)
  formatterMock.FormatMock.Return("minimock")

  //tested code can run the mocked method in a goroutine
  go formatterMock.Format("hello world!")
}
```

## Using GoUnit with minimock

Writing test is not only mocking the dependencies. Often the test itself contains a lot of boilerplate code.
You can generate test stubs using [GoUnit](https://github.com/hexdigest/gounit) tool which has a nice template that uses minimock.

Happy mocking!
