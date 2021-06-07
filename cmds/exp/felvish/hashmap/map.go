package hashmap

// Map is a persistent associative data structure mapping keys to values. It
// is immutable, and supports near-O(1) operations to create modified version of
// the map that shares the underlying data structure. Because it is immutable,
// all of its methods are safe for concurrent use.
type Map interface {
	Len() int
	// Index returns whether there is a value associated with the given key, and
	// that value or nil.
	Index(k interface{}) (interface{}, bool)
	// Assoc returns an almost identical map, with the given key associated with
	// the given value.
	Assoc(k, v interface{}) Map
	// Dissoc returns an almost identical map, with the given key associated
	// with no value.
	Dissoc(k interface{}) Map
	// Iterator returns an iterator over the map.
	Iterator() Iterator
}

// Iterator is an iterator over map elements. It can be used like this:
//
//     for it := m.Iterator(); it.HasElem(); it.Next() {
//         key, value := it.Elem()
//         // do something with elem...
//     }
type Iterator interface {
	// Elem returns the current key-value pair.
	Elem() (interface{}, interface{})
	// HasElem returns whether the iterator is pointing to an element.
	HasElem() bool
	// Next moves the iterator to the next position.
	Next()
}

// HasKey reports whether a Map has the given key.
func HasKey(m Map, k interface{}) bool {
	_, ok := m.Index(k)
	return ok
}
