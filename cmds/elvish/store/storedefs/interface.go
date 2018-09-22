package storedefs

// Store is an interface for things like command history
// It assumes we have a 'cursor', which we can adjust.
type Store interface {
	Seek(int) (int, error)
	Cur() (int, error)
	Prev() (int, error)
	Next() (int, error)
	Remove(int) error
	Add(text string) (int, error)
	One(int) (string, error)
	List(int, int) ([]string, error)
	Walk(from int, to int, f func(string) bool) error
	Search(from int, prefix string) (int, string, error)
	RSearch(upto int, prefix string) (int, string, error)
}
