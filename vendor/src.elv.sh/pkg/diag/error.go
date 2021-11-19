package diag

import "fmt"

// Error represents an error with context that can be showed.
type Error struct {
	Type    string
	Message string
	Context Context
}

// Error returns a plain text representation of the error.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %d-%d in %s: %s",
		e.Type, e.Context.From, e.Context.To, e.Context.Name, e.Message)
}

// Range returns the range of the error.
func (e *Error) Range() Ranging {
	return e.Context.Range()
}

// Show shows the error.
func (e *Error) Show(indent string) string {
	header := fmt.Sprintf("%s: \033[31;1m%s\033[m\n", e.Type, e.Message)
	return header + e.Context.ShowCompact(indent+"  ")
}
