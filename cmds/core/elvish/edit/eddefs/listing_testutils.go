package eddefs

import (
	"fmt"
	"reflect"

	"github.com/u-root/u-root/cmds/core/elvish/edit/ui"
)

type ListingShown struct {
	Header  string
	Content ui.Styled
}

type ListingProviderFilterTest struct {
	Filter     string
	WantShowns []ListingShown
}

func TestListingProviderFilter(name string, ls ListingProvider, testcases []ListingProviderFilterTest) error {
	for _, testcase := range testcases {
		ls.Filter(testcase.Filter)

		l := ls.Len()
		if l != len(testcase.WantShowns) {
			return fmt.Errorf("%s.Len() -> %d, want %d (filter was %q)",
				name, l, len(testcase.WantShowns), testcase.Filter)
		} else {
			for i, want := range testcase.WantShowns {
				header, content := ls.Show(i)
				if header != want.Header || !reflect.DeepEqual(content, want.Content) {
					return fmt.Errorf("%s.Show(%d) => (%v, %v), want (%v, %v) (filter was %q)",
						name, i, header, content, want.Header, want.Content, testcase.Filter)
				}
			}
		}
	}
	return nil
}
