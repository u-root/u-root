package tc

import "fmt"

func concatError(existing, new error) error {
	if new == nil {
		return existing
	}
	if existing == nil {
		return new
	}
	return fmt.Errorf("%v\n%v", existing, new)
}
