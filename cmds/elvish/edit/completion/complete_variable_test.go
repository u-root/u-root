package completion

import (
	"reflect"
	"sort"
	"testing"

	"github.com/u-root/u-root/cmds/elvish/parse"
)

func TestFindVariableComplContext(t *testing.T) {
	testComplContextFinder(t, "findVariableComplContext", findVariableComplContext, []complContextFinderTest{
		{"$", &variableComplContext{
			complContextCommon{"", parse.Bareword, 1, 1}, "", ""}},
		{"$a", &variableComplContext{
			complContextCommon{"a", parse.Bareword, 1, 2}, "", ""}},
		{"$a:", &variableComplContext{
			complContextCommon{"", parse.Bareword, 3, 3}, "a", "a:"}},
		{"$a:b", &variableComplContext{
			complContextCommon{"b", parse.Bareword, 3, 4}, "a", "a:"}},
		// Wrong contexts
		{"", nil},
		{"echo", nil},
	})
}

type testEvalerScopes struct{}

var testScopes = map[string]map[string]int{
	"":        {"veni": 0, "vidi": 0, "vici": 0},
	"foo":     {"lorem": 0, "ipsum": 0},
	"foo:bar": {"lorem": 0, "dolor": 0},
}

func (testEvalerScopes) EachNsInTop(f func(string)) {
	for ns := range testScopes {
		if ns != "" {
			f(ns)
		}
	}
}

func (testEvalerScopes) EachVariableInTop(ns string, f func(string)) {
	for name := range testScopes[ns] {
		f(name)
	}
}

var complVariableTests = []struct {
	ns     string
	nsPart string
	want   []rawCandidate
}{
	// No namespace: complete variables and namespaces
	{"", "", []rawCandidate{
		noQuoteCandidate("foo:"), noQuoteCandidate("foo:bar:"),
		noQuoteCandidate("veni"), noQuoteCandidate("vici"), noQuoteCandidate("vidi"),
	}},
	// Nonempty namespace: complete variables in namespace and subnamespaces
	// (but not variables in subnamespaces)
	{"foo", "foo:", []rawCandidate{
		noQuoteCandidate("bar:"),
		noQuoteCandidate("ipsum"), noQuoteCandidate("lorem"),
	}},
	// Bad namespace
	{"bad", "bad:", nil},
}

func TestComplVariable(t *testing.T) {
	for _, test := range complVariableTests {
		got := collectComplVariable(test.ns, test.nsPart, testEvalerScopes{})
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("complVariable(%q, %q, ...) => %v, want %v", test.ns, test.nsPart, got, test.want)
		}
	}
}

func collectComplVariable(ns, nsPart string, ev evalerScopes) []rawCandidate {
	ch := make(chan rawCandidate)
	go func() {
		complVariable(ns, nsPart, ev, ch)
		close(ch)
	}()
	var results []rawCandidate
	for result := range ch {
		results = append(results, result)
	}
	sort.Sort(rawCandidates(results))
	return results
}
