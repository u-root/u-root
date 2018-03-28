package complete

import (
	"strings"
)

type StringCompleter struct {
	Names []string
}

func NewStringCompleter(s []string) Completer {
	return &StringCompleter{Names: s}
}

func (f *StringCompleter) Complete(s string) ([]string, error) {
	var names []string
	for _, n := range f.Names {
		Debug("Check %v against %v", n, s)
		if strings.HasPrefix(n, s) {
			Debug("Add %v", n)
			names = append(names, n)
		}
	}
	return names, nil
}
