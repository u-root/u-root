package history

import (
	"testing"

	"github.com/u-root/u-root/cmds/core/elvish/edit/eddefs"
	"github.com/u-root/u-root/cmds/core/elvish/edit/ui"
)

var (
	theHistList = newHistlist([]string{"ls", "echo lalala", "ls"})

	histlistDedupFilterTests = []eddefs.ListingProviderFilterTest{
		{"", []eddefs.ListingShown{
			{"1", ui.Unstyled("echo lalala")},
			{"2", ui.Unstyled("ls")}}},
		{"l", []eddefs.ListingShown{
			{"1", ui.Unstyled("echo lalala")},
			{"2", ui.Unstyled("ls")}}},
	}

	histlistNoDedupFilterTests = []eddefs.ListingProviderFilterTest{
		{"", []eddefs.ListingShown{
			{"0", ui.Unstyled("ls")},
			{"1", ui.Unstyled("echo lalala")},
			{"2", ui.Unstyled("ls")}}},
		{"l", []eddefs.ListingShown{
			{"0", ui.Unstyled("ls")},
			{"1", ui.Unstyled("echo lalala")},
			{"2", ui.Unstyled("ls")}}},
	}
)

func TestHistlist(t *testing.T) {
	if err := eddefs.TestListingProviderFilter("theHistList", theHistList, histlistDedupFilterTests); err != nil {
		t.Error(err)
	}
	theHistList.dedup = false
	if err := eddefs.TestListingProviderFilter("theHistList", theHistList, histlistNoDedupFilterTests); err != nil {
		t.Error(err)
	}
}
