package complete

import (
	"fmt"
	"reflect"
	"sort"
)

// StringEntry adds the Entry interface to a simple string. Its description part is empty.
func StringEntry(s string) Entry { return stringEntry(s) }

type stringEntry string

func (s stringEntry) Title() string { return string(s) }

func (stringEntry) Description() string { return "" }

// StringValues adds the Values interface to a simple string
// slice. There is just one category.
func StringValues(title string, entries []string) Values {
	return stringValues{title, entries}
}

type stringValues struct {
	title   string
	entries []string
}

func (s stringValues) NumCategories() int              { return 1 }
func (s stringValues) CategoryTitle(_ int) string      { return s.title }
func (s stringValues) NumEntries(_ int) int            { return len(s.entries) }
func (s stringValues) Entry(_ int, entryIdx int) Entry { return StringEntry(s.entries[entryIdx]) }

// MapValues adds the Values interface to a map of entries.
//
// In go 1.18, this function would be:
//
//    func MapValues[T Entry](values map[string][]T, categories []string)
//
// Each of the map values should be a slice of objects implementing
// the Entry interface.
//
// The categories string slice, if provided, selects a specific order
// for the categories. If nil is specified, the map keys are used
// in sorted order.
func MapValues(values interface{}, categories []string) Values {
	m := reflect.ValueOf(values)
	if categories == nil {
		categories = make([]string, 0, m.Len())
		for _, k := range m.MapKeys() {
			categories = append(categories, k.Interface().(string))
		}
		sort.Strings(categories)
	}
	return mapValues{m, categories}
}

type mapValues struct {
	values     reflect.Value
	categories []string
}

func (s mapValues) NumCategories() int         { return len(s.categories) }
func (s mapValues) CategoryTitle(i int) string { return s.categories[i] }
func (s mapValues) NumEntries(i int) int {
	entries := s.values.MapIndex(reflect.ValueOf(s.categories[i]))
	return entries.Len()
}
func (s mapValues) Entry(cat, entry int) Entry {
	slice := s.values.MapIndex(reflect.ValueOf(s.categories[cat]))
	val := slice.Index(entry)
	e, ok := val.Interface().(Entry)
	if !ok {
		e, ok = val.Addr().Interface().(Entry)
		if !ok {
			if val.Type().Kind() == reflect.String {
				e = StringEntry(val.Interface().(string))
			} else {
				panic(fmt.Errorf("can't convert %v to Entry", val.Type()))
			}
		}
	}
	return e
}
