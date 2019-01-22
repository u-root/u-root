// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complete

// MultiCompleter is a Completer consisting of one or more Completers
// Why do this?
// We need it for paths, anyway, but consider a shell which
// has builtins and metacharacters such as >, &, etc.
// You can build a MultiCompleter which has a string completer
// and a set of file completers, so you don't need to special
// case anything.
type MultiCompleter struct {
	Completers []Completer
}

// NewMultiCompleter returns a MultiCompleter created from
// one or more Completers. It is perfectly legal to include a
// MultiCompleter.
func NewMultiCompleter(c Completer, cc ...Completer) Completer {
	return &MultiCompleter{append([]Completer{c}, cc...)}
}

// Complete Returns a []string consisting of the results
// of calling all the Completers.
func (m *MultiCompleter) Complete(s string) (string, []string, error) {
	var files []string
	var exact string
	for _, c := range m.Completers {
		x, cc, err := c.Complete(s)
		if err != nil {
			Debug("MultiCompleter: %v: %v", c, err)
		}
		files = append(files, cc...)
		if exact == "" {
			exact = x
		}

	}
	return exact, files, nil
}
