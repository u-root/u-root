![logo](https://rawgit.com/gojuno/minimock/master/logo.svg)
[![GoDoc](https://godoc.org/github.com/gojuno/minimock?status.svg)](http://godoc.org/github.com/gojuno/minimock) [![Build Status](https://travis-ci.org/gojuno/minimock.svg?branch=master)](https://travis-ci.org/gojuno/minimock) [![Go Report Card](https://goreportcard.com/badge/github.com/gojuno/minimock)](https://goreportcard.com/report/github.com/gojuno/minimock) [![cover.run](https://cover.run/go/github.com/gojuno/minimock.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2Fgojuno%2Fminimock) [![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go#testing)


## Summary 
Minimock parses the input Go source file that contains an interface declaration and generates an
implementation of this interface that can be used as a mock.

Main features of minimock:

* It's integrated with the standard Go "testing" package
* It supports variadic methods and embedded interfaces
* It's very convenient to use generated mocks in table tests because it implements builder pattern to set up several mocks
* It provides a useful Controller.Wait(time.Duration) helper method to test concurrent code
* It generates helpers to check if the mocked methods have been called and keeps your tests clean and up to date
* It generates concurrent-safe mock execution counters that you can use in your mocks to implement sophisticated mocks behavior
* There is the template for [GoUnit](https://github.com/hexdigest/gounit) tool which generates table driven tests that make use of minimock

## Installation

```
go get github.com/gojuno/minimock/cmd/minimock
```

## Usage
Let's say we have the following interface declaration in github.com/gojuno/minimock/tests package:
```go
type Formatter interface {
	Format(string, ...interface{}) string
}
```

Here is how to generate the mock for this interface:
```
minimock -i github.com/gojuno/minimock/tests.Formatter -o ./tests/
```

You can skip the package name if you are in the directory where `Formater` is defined:
```
cd ~/go/src/github.com/gojuno/minimock/tests
minimock -i Formatter -o ./tests/
```

The result file ./tests/formatter_mock_test.go will be:
```go

package tests

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.8
The original interface "Formatter" can be found in github.com/gojuno/minimock/tests
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	testify_assert "github.com/stretchr/testify/assert"
)

//FormatterMock implements github.com/gojuno/minimock/tests.Formatter
type FormatterMock struct {
	t minimock.Tester

	FormatFunc    func(p string, p1 ...interface{}) (r string)
	FormatCounter uint64
	FormatMock    mFormatterMockFormat
}

//NewFormatterMock returns a mock for github.com/gojuno/minimock/tests.Formatter
func NewFormatterMock(t minimock.Tester) *FormatterMock {
	m := &FormatterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.FormatMock = mFormatterMockFormat{mock: m}

	return m
}

type mFormatterMockFormat struct {
	mock             *FormatterMock
	mockExpectations *FormatterMockFormatParams
}

//FormatterMockFormatParams represents input parameters of the Formatter.Format
type FormatterMockFormatParams struct {
	p  string
	p1 []interface{}
}

//Expect sets up expected params for the Formatter.Format
func (m *mFormatterMockFormat) Expect(p string, p1 ...interface{}) *mFormatterMockFormat {
	m.mockExpectations = &FormatterMockFormatParams{p, p1}
	return m
}

//Return sets up a mock for Formatter.Format to return Return's arguments
func (m *mFormatterMockFormat) Return(r string) *FormatterMock {
	m.mock.FormatFunc = func(p string, p1 ...interface{}) string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Formatter.Format method
func (m *mFormatterMockFormat) Set(f func(p string, p1 ...interface{}) (r string)) *FormatterMock {
	m.mock.FormatFunc = f
	return m.mock
}

//Format implements github.com/gojuno/minimock/tests.Formatter interface
func (m *FormatterMock) Format(p string, p1 ...interface{}) (r string) {
	defer atomic.AddUint64(&m.FormatCounter, 1)

	if m.FormatMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.FormatMock.mockExpectations, FormatterMockFormatParams{p, p1},
			"Formatter.Format got unexpected parameters")

		if m.FormatFunc == nil {
			m.t.Fatal("No results are set for the FormatterMock.Format")
			return
		}
	}

	if m.FormatFunc == nil {
		m.t.Fatal("Unexpected call to FormatterMock.Format")
		return
	}

	return m.FormatFunc(p, p1...)
}

//FormatMinimockCounter returns a count of Formatter.Format invocations
func (m *FormatterMock) FormatMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FormatCounter)
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *FormatterMock) MinimockFinish() {
	if m.FormatFunc != nil && atomic.LoadUint64(&m.FormatCounter) == 0 {
		m.t.Fatal("Expected call to FormatterMock.Format")
	}
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *FormatterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.FormatFunc == nil || atomic.LoadUint64(&m.FormatCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.FormatFunc != nil && atomic.LoadUint64(&m.FormatCounter) == 0 {
				m.t.Error("Expected call to FormatterMock.Format")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
```


There are several ways to set up a mock

Setting up a mock using direct assignment:
```go
formatterMock := NewFormatterMock(mc)
formatterMock.FormatFunc = func(string, ...interface{}) string {
  return "minimock"
}
```

Setting up a mock using builder pattern and Expect/Return methods:
```go
formatterMock := NewFormatterMock(mc).FormatMock.Expect("%s %d", "string", 1).Return("minimock")
```

Setting up a mock using builder and Set method:
```go
formatterMock := NewFormatterMock(mc).FormatMock.Set(func(string, ...interface{}) string {
  return "minimock"
})
```

Builder pattern is convenient when you have to mock more than one method of an interface.
Let's say we have an io.ReadCloser interface which has two methods: Read and Close
```go
type ReadCloser interface {
	Read(p []byte) (n int, err error)
	Close() error
}
```

Then you can set up a mock using just one assignment:
```go
readCloserMock := NewReadCloserMock(mc).ReadMock.Expect([]byte(1,2,3)).Return(3, nil).CloseMock.Return(nil)
```

You can also use invocation counters in your mocks and tests:
```go
formatterMock := NewFormatterMock(mc)
formatterMock.FormatFunc = func(string, ...interface{}) string {
  return fmt.Sprintf("minimock: %d", formatterMock.FormatMinimockCounter())
}
```

## minimock.Controller
When you have to mock multiple dependencies in your test it's recommended to use minimock.Controller and its Finish or Wait methods.
All you have to do is instantiate the Controller and pass it as an argument to the mocks' constructors:
```go
func TestSomething(t *testing.T) {
  mc := minimock.NewController(t)
  defer mc.Finish()

  formatterMock := NewFormatterMock(mc)
  formatterMock.FormatMock.Return("minimock")

  readCloserMock := NewReadCloserMock(mc)
  readCloserMock.ReadMock.Return(5, nil)

  readCloserMock.Read([]byte{}) 
  formatterMock.Format()
}
```
Every mock is registered in the controller so by calling mc.Finish() you can verify that all the registered mocks have been called
within your test.

## Keep your tests clean
Sometimes we write tons of mocks for our tests but over time the tested code stops using mocked dependencies,
however mocks are still present and being initialized in the test files. So while tested code can shrink, tests are only growing.
To prevent this minimock.Controller provides Finish() method that verifies that all your mocks have been called at least once during the test run.

```go
func TestSomething(t *testing.T) {
  mc := minimock.NewController(t)
  defer mc.Finish() //this will mark your test as failed because there are no calls to formatterMock.Format() and readCloserMock.Read() below

  formatterMock := NewFormatterMock(mc)
  formatterMock.FormatMock.Return("minimock")

  readCloserMock := NewReadCloserMock(mc)
  readCloserMock.ReadMock.Return(5, nil)
}
```

## Testing concurrent code
Testing concurrent code is tough. Fortunately minimock provides you with the helper method that makes testing concurrent code easy.
Here is how it works:

```go
func TestSomething(t *testing.T) {
  mc := minimock.NewController(t)

  //Wait ensures that all mocked methods have been called within given interval
  //if any of the mocked methods have not been called Wait marks test as failed
  defer mc.Wait(time.Second)

  formatterMock := NewFormatterMock(mc)
  formatterMock.FormatMock.Return("minimock")

  //tested code can run mocked method in a goroutine
  go formatterMock.Format("")
}
```

## Using gotests with minimock

Despite the fact that minimock does a lot of work for you, writing test is not only mocking the dependencies.
There is a nice tool called [GoTests](https://github.com/hexdigest/gotests) that generates table driven tests with minimock in mind.

Let's say you have to test the following code:
```go
type Dependency interface {
	SomeMethod() error
}

type Service struct {
	Dep Dependency
}

func (s Service) DoSomething() bool {
	if err := s.Dep.SomeMethod(); err != nil {
		log.Printf("some method returned an error: %v", err)
		return false
	}

	return true
}
```

GoTests generates the following test for the DoSomething method:
```go
func TestService_DoSomething(t *testing.T) {
	tests := []struct {
		name  string
		setup func(mc *minimock.Controller) Service
		want  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)
			s := tt.setup(mc)
			if got := s.DoSomething(); got != tt.want {
				t.Errorf("Service.DoSomething() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

So all you need here is to fill the tests slice with:
```go
{
  name: "some method failed",
  setup: func(mc *minimock.Controller) Service {
    return Service{
      Dep: NewDependencyMock(mc).SomeMethodMock.Return(io.EOF),
    }
  },
  want: false,
},
{
  name: "some method succeeded",
  setup: func(mc *minimock.Controller) Service {
    return Service{
      Dep: NewDependencyMock(mc).SomeMethodMock.Return(nil),
    }
  },
  want: true,
},
```

And your test is ready!

## Minimock command line flags:
```
$ minimock -h
  Usage of minimock:
    -f string
      	DEPRECATED: input file or import path of the package that contains interface declaration
    -h	show this help message
    -i string
      	comma-separated names of the interfaces to mock, i.e fmt.Stringer,io.Reader, use io.* notation to generate mocks for all interfaces in an io package
    -o string
      	destination file name to place the generated mock or path to destination package when multiple interfaces are given
    -p string
      	DEPRECATED: destination package name
    -s string
      	output file name suffix which is added to file names when multiple interfaces are given (default "_mock_test.go")
    -t string
      	DEPRECATED: mock struct name (default <interface name>Mock)
    -withTests
      	parse *_test.go files in the source package
```
