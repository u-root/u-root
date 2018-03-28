package complete

type Completer interface {
	Complete(s string) ([]string, error)
}

var Debug = func(s string, v ...interface{}) {}
