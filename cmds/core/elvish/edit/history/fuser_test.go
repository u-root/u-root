package history

import (
	"testing"

	"github.com/u-root/u-root/cmds/core/elvish/store"
)

var fuserStore = store.NewCmdHistory("store 1")

func TestFuser(t *testing.T) {
	_, err := NewFuser(fuserStore)
	if err != nil {
		t.Errorf("NewFuser -> error %v, want nil", err)
	}
	// fix this fucking mess later.

}
