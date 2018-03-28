package complete

type Completer interface {
	Complete(s string) ([]string, error)
}

var debug = func(s string, v ...interface{}) {}
