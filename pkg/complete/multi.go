package complete

import (
	"fmt"
	"os"
	"strings"
)

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

func NewMultiCompleter(c Completer, cc ...Completer) Completer {
	return &MultiCompleter{append([]Completer{c}, cc...)}
}

// Complete Returns a []string consisting of the results
// of calling all the Completers.
func (m *MultiCompleter) Complete(s string) ([]string, error) {
	var files []string
	for _, c := range m.Completers {
		cc, err := c.Complete(s)
		if err != nil {
			Debug("MultiCompleter: %v: %v", c, err)
		}
		files = append(files, cc...)
	}
	return files, nil
}

func NewEnvCompleter(s string) (Completer, error) {
	dirs := strings.Split(s, ":")
	if len(dirs) == 0 {
		return nil, fmt.Errorf("%s is empty", s)
	}
	c := make([]Completer, len(dirs))
	for i := range dirs {
		c[i] = NewFileCompleter(dirs[i])
	}
	return NewMultiCompleter(c[0], c[1:]...), nil
}

func NewPathCompleter() (Completer, error) {
	// Getenv returns the same value ("") if a path is not found
	// or if it has the value "". Oh well.
	return NewEnvCompleter(os.Getenv("PATH"))
}
