package edit

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"src.elv.sh/pkg/cli/mode"
	"src.elv.sh/pkg/cli/tk"
	"src.elv.sh/pkg/eval"
	"src.elv.sh/pkg/eval/vals"
	"src.elv.sh/pkg/eval/vars"
	"src.elv.sh/pkg/strutil"
	"src.elv.sh/pkg/ui"
)

type customListingOpts struct {
	Binding    bindingsMap
	Caption    string
	KeepBottom bool
	Accept     eval.Callable
	AutoAccept bool
}

func (*customListingOpts) SetDefaultOptions() {}

//elvdoc:fn listing:start-custom
//
// Starts a custom listing addon.

func listingStartCustom(ed *Editor, fm *eval.Frame, opts customListingOpts, items interface{}) {
	var bindings tk.Bindings
	if opts.Binding.Map != nil {
		bindings = newMapBindings(ed, fm.Evaler, vars.FromPtr(&opts.Binding))
	}
	var getItems func(string) []mode.ListingItem
	if fn, isFn := items.(eval.Callable); isFn {
		getItems = func(q string) []mode.ListingItem {
			var items []mode.ListingItem
			var itemsMutex sync.Mutex
			collect := func(item mode.ListingItem) {
				itemsMutex.Lock()
				defer itemsMutex.Unlock()
				items = append(items, item)
			}
			valuesCb := func(ch <-chan interface{}) {
				for v := range ch {
					if item, itemOk := getListingItem(v); itemOk {
						collect(item)
					}
				}
			}
			bytesCb := func(r *os.File) {
				buffered := bufio.NewReader(r)
				for {
					line, err := buffered.ReadString('\n')
					if line != "" {
						s := strutil.ChopLineEnding(line)
						collect(mode.ListingItem{ToAccept: s, ToShow: ui.T(s)})
					}
					if err != nil {
						break
					}
				}
			}
			f := func(fm *eval.Frame) error { return fn.Call(fm, []interface{}{q}, eval.NoOpts) }
			err := fm.PipeOutput(f, valuesCb, bytesCb)
			// TODO(xiaq): Report the error.
			_ = err
			return items
		}
	} else {
		getItems = func(q string) []mode.ListingItem {
			convertedItems := []mode.ListingItem{}
			vals.Iterate(items, func(v interface{}) bool {
				toFilter, toFilterOk := getToFilter(v)
				item, itemOk := getListingItem(v)
				if toFilterOk && itemOk && strings.Contains(toFilter, q) {
					// TODO(xiaq): Report type error when ok is false.
					convertedItems = append(convertedItems, item)
				}
				return true
			})
			return convertedItems
		}
	}

	w, err := mode.NewListing(ed.app, mode.ListingSpec{
		Bindings: bindings,
		Caption:  opts.Caption,
		GetItems: func(q string) ([]mode.ListingItem, int) {
			items := getItems(q)
			selected := 0
			if opts.KeepBottom {
				selected = len(items) - 1
			}
			return items, selected
		},
		Accept: func(s string) bool {
			if opts.Accept != nil {
				callWithNotifyPorts(ed, fm.Evaler, opts.Accept, s)
			}
			return false
		},
		AutoAccept: opts.AutoAccept,
	})
	startMode(ed.app, w, err)
}

func getToFilter(v interface{}) (string, bool) {
	toFilterValue, _ := vals.Index(v, "to-filter")
	toFilter, toFilterOk := toFilterValue.(string)
	return toFilter, toFilterOk
}

func getListingItem(v interface{}) (item mode.ListingItem, ok bool) {
	toAcceptValue, _ := vals.Index(v, "to-accept")
	toAccept, toAcceptOk := toAcceptValue.(string)
	toShowValue, _ := vals.Index(v, "to-show")
	toShow, toShowOk := toShowValue.(ui.Text)
	if toShowString, ok := toShowValue.(string); ok {
		toShow = ui.T(toShowString)
		toShowOk = true
	}
	return mode.ListingItem{ToAccept: toAccept, ToShow: toShow}, toAcceptOk && toShowOk
}
