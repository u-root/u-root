package lastcmd

import (
	"testing"

	"github.com/u-root/u-root/cmds/core/elvish/edit/eddefs"
	"github.com/u-root/u-root/cmds/core/elvish/edit/ui"
)

var (
	theLine    = "qw search 'foo bar ~y'"
	theLastCmd = newState(theLine)

	tests = []eddefs.ListingProviderFilterTest{
		{"", []eddefs.ListingShown{
			{"M-1", ui.Unstyled(theLine)},
			{"0", ui.Unstyled("qw")},
			{"1", ui.Unstyled("search")},
			{"2", ui.Unstyled("'foo bar ~y'")}}},
		{"1", []eddefs.ListingShown{{"1", ui.Unstyled("search")}}},
		{"-", []eddefs.ListingShown{
			{"M-1", ui.Unstyled(theLine)},
			{"-3", ui.Unstyled("qw")},
			{"-2", ui.Unstyled("search")},
			{"-1", ui.Unstyled("'foo bar ~y'")}}},
		{"-1", []eddefs.ListingShown{{"-1", ui.Unstyled("'foo bar ~y'")}}},
	}
)

func TestLastCmd(t *testing.T) {
	if err := eddefs.TestListingProviderFilter("theLastCmd", theLastCmd, tests); err != nil {
		t.Error(err)
	}
}
