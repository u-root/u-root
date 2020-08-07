package tests

import (
	"reflect"
	"testing"
	"time"

	minimock "github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatterMock_ImplementsStringer(t *testing.T) {
	v := NewFormatterMock(NewTesterMock(t))
	assert.True(t, reflect.TypeOf(v).Implements(reflect.TypeOf((*Formatter)(nil)).Elem()))
}

func TestFormatterMock_UnmockedCallFailsTest(t *testing.T) {
	var mockCalled bool
	tester := NewTesterMock(t)
	tester.FatalfMock.Set(func(s string, args ...interface{}) {
		assert.Equal(t, "Unexpected call to FormatterMock.Format. %v %v", s)
		assert.Equal(t, "this call fails because Format method isn't mocked", args[0])

		mockCalled = true
	})

	defer tester.MinimockFinish()

	formatterMock := NewFormatterMock(tester)
	dummyFormatter{formatterMock}.Format("this call fails because Format method isn't mocked")
	assert.True(t, mockCalled)
}

func TestFormatterMock_MockedCallSucceeds(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewFormatterMock(tester)
	formatterMock.FormatMock.Set(func(format string, args ...interface{}) string {
		return "mock is successfully called"
	})
	defer tester.MinimockFinish()

	df := dummyFormatter{formatterMock}
	assert.Equal(t, "mock is successfully called", df.Format(""))
}

func TestFormatterMock_Wait(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewFormatterMock(tester)
	formatterMock.FormatMock.Set(func(format string, args ...interface{}) string {
		return "mock is successfully called from the goroutine"
	})

	go func() {
		df := dummyFormatter{formatterMock}
		assert.Equal(t, "mock is successfully called from the goroutine", df.Format(""))
	}()

	formatterMock.MinimockWait(time.Second)
}

func TestFormatterMock_Expect(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewFormatterMock(tester).FormatMock.Expect("Hello", "world", "!").Return("")

	df := dummyFormatter{formatterMock}
	df.Format("Hello", "world", "!")

	assert.EqualValues(t, 1, formatterMock.FormatBeforeCounter())
	assert.EqualValues(t, 1, formatterMock.FormatAfterCounter())
}

func TestFormatterMock_ExpectDifferentArguments(t *testing.T) {
	assert.Panics(t, func() {
		tester := NewTesterMock(t)
		defer tester.MinimockFinish()

		tester.ErrorfMock.Set(func(s string, args ...interface{}) {
			assert.Equal(t, "FormatterMock.Format got unexpected parameters, want: %#v, got: %#v%s\n", s)
			require.Len(t, args, 3)
			assert.Equal(t, FormatterMockFormatParams{s1: "expected"}, args[0])
			assert.Equal(t, FormatterMockFormatParams{s1: "actual"}, args[1])
		})

		tester.FatalMock.Expect("No results are set for the FormatterMock.Format").Return()

		formatterMock := NewFormatterMock(tester)
		formatterMock.FormatMock.Expect("expected")
		formatterMock.Format("actual")
	})
}

func TestFormatterMock_ExpectAfterSet(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.FatalfMock.Expect("FormatterMock.Format mock is already set by Set").Return()

	formatterMock := NewFormatterMock(tester)
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })

	formatterMock.FormatMock.Expect("Should not work")
}

func TestFormatterMock_ExpectAfterWhen(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.FatalfMock.Expect("Expectation set by When has same params: %#v", FormatterMockFormatParams{s1: "Should not work", p1: nil}).Return()

	formatterMock := NewFormatterMock(tester)
	formatterMock.FormatMock.When("Should not work").Then("")

	formatterMock.Format("Should not work")

	formatterMock.FormatMock.Expect("Should not work")
}

func TestFormatterMock_Return(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewFormatterMock(tester).FormatMock.Return("Hello world!")
	df := dummyFormatter{formatterMock}
	assert.Equal(t, "Hello world!", df.Format(""))
}

func TestFormatterMock_ReturnAfterSet(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.FatalfMock.Expect("FormatterMock.Format mock is already set by Set").Return()

	formatterMock := NewFormatterMock(tester)
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })

	formatterMock.FormatMock.Return("Should not work")
}

func TestFormatterMock_ReturnWithoutExpectForFixedArgsMethod(t *testing.T) {
	// Test for issue https://github.com/gojuno/minimock/issues/31

	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.ErrorMock.Expect("Expected call to FormatterMock.Format")
	tester.FailNowMock.Expect()

	formatterMock := NewFormatterMock(tester)
	formatterMock.FormatMock.Return("")
	formatterMock.MinimockFinish()
}

func TestFormatterMock_Set(t *testing.T) {
	tester := NewTesterMock(t)

	formatterMock := NewFormatterMock(tester).FormatMock.Set(func(string, ...interface{}) string {
		return "set"
	})

	df := dummyFormatter{formatterMock}
	assert.Equal(t, "set", df.Format(""))
}

func TestFormatterMock_SetAfterExpect(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.FatalfMock.Expect("Default expectation is already set for the Formatter.Format method").Return()

	formatterMock := NewFormatterMock(tester).FormatMock.Expect("").Return("")

	//second attempt should fail
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })
}

func TestFormatterMock_SetAfterWhen(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.FatalfMock.Expect("Some expectations are already set for the Formatter.Format method").Return()

	formatterMock := NewFormatterMock(tester).FormatMock.When("").Then("")

	//second attempt should fail
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })
}

func TestFormatterMockFormat_WhenThen(t *testing.T) {
	formatter := NewFormatterMock(t)
	defer formatter.MinimockFinish()

	formatter.FormatMock.When("hello %v", "username").Then("hello username")
	formatter.FormatMock.When("goodbye %v", "username").Then("goodbye username")

	assert.Equal(t, "hello username", formatter.Format("hello %v", "username"))
	assert.Equal(t, "goodbye username", formatter.Format("goodbye %v", "username"))
}

func TestFormatterMockFormat_WhenAfterSet(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.FatalfMock.Expect("FormatterMock.Format mock is already set by Set").Return()

	formatterMock := NewFormatterMock(tester)
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })

	formatterMock.FormatMock.When("Should not work")
}

func TestFormatterMock_MinimockFormatDone(t *testing.T) {
	formatterMock := NewFormatterMock(t)

	formatterMock.FormatMock.expectations = []*FormatterMockFormatExpectation{{}}
	assert.False(t, formatterMock.MinimockFormatDone())

	formatterMock = NewFormatterMock(t)
	formatterMock.FormatMock.defaultExpectation = &FormatterMockFormatExpectation{}
	assert.False(t, formatterMock.MinimockFormatDone())
}

func TestFormatterMock_MinimockFinish(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.ErrorMock.Expect("Expected call to FormatterMock.Format").Return()
	tester.FailNowMock.Expect().Return()

	formatterMock := NewFormatterMock(tester)
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })

	formatterMock.MinimockFinish()
}

func TestFormatterMock_MinimockFinish_WithNoMetExpectations(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.ErrorfMock.Set(func(m string, args ...interface{}) {
		assert.Equal(t, m, "Expected call to FormatterMock.Format with params: %#v")
	})
	tester.FailNowMock.Expect().Return()

	formatterMock := NewFormatterMock(tester)
	formatterMock.FormatMock.Expect("a").Return("a")
	formatterMock.FormatMock.When("b").Then("b")

	formatterMock.MinimockFinish()
}

func TestFormatterMock_MinimockWait(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	tester.ErrorMock.Expect("Expected call to FormatterMock.Format").Return()
	tester.FailNowMock.Expect().Return()

	formatterMock := NewFormatterMock(tester)
	formatterMock.FormatMock.Set(func(string, ...interface{}) string { return "" })

	formatterMock.MinimockWait(time.Millisecond)
}

// Verifies that Calls() doesn't return nil if no calls were made
func TestFormatterMock_CallsNotNil(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	formatterMock := NewFormatterMock(tester)
	calls := formatterMock.FormatMock.Calls()

	assert.NotNil(t, calls)
	assert.Empty(t, calls)
}

// Verifies that Calls() returns the correct call args in the expected order
func TestFormatterMock_Calls(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	// Arguments used for each mock call
	expected := []*FormatterMockFormatParams{
		{"a1", []interface{}{}},
		{"b1", []interface{}{"b2"}},
		{"c1", []interface{}{"c2", "c3"}},
		{"d1", []interface{}{"d2", "d3", "d4"}},
	}

	formatterMock := NewFormatterMock(tester)

	for _, p := range expected {
		formatterMock.FormatMock.Expect(p.s1, p.p1...).Return("")
		formatterMock.Format(p.s1, p.p1...)
	}

	assert.Equal(t, expected, formatterMock.FormatMock.Calls())
}

// Verifies that Calls() returns a new shallow copy of the params list each time
func TestFormatterMock_CallsReturnsCopy(t *testing.T) {
	tester := NewTesterMock(t)
	defer tester.MinimockFinish()

	expected := []*FormatterMockFormatParams{
		{"a1", []interface{}{"a1"}},
		{"b1", []interface{}{"b2"}},
	}

	formatterMock := NewFormatterMock(tester)
	callHistory := [][]*FormatterMockFormatParams{}

	for _, p := range expected {
		formatterMock.FormatMock.Expect(p.s1, p.p1...).Return("")
		formatterMock.Format(p.s1, p.p1...)
		callHistory = append(callHistory, formatterMock.FormatMock.Calls())
	}

	assert.Equal(t, len(expected), len(callHistory))

	for i, c := range callHistory {
		assert.Equal(t, i+1, len(c))
	}
}

type dummyFormatter struct {
	Formatter
}

type dummyMockController struct {
	minimock.MockController
	registerCounter int
}

func (dmc *dummyMockController) RegisterMocker(m minimock.Mocker) {
	dmc.registerCounter++
}

func TestFormatterMock_RegistersMocker(t *testing.T) {
	mockController := &dummyMockController{}

	NewFormatterMock(mockController)
	assert.Equal(t, 1, mockController.registerCounter)
}
