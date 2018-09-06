package uflag

import "strings"

// Strings implements flag.Value that appends multiple invocations of the
// flag to a slice of strings.
type Strings []string

// Set implements flag.Value.Set.
func (s *Strings) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// String implements flag.Value.String.
func (s Strings) String() string {
	return strings.Join(s, ",")
}
