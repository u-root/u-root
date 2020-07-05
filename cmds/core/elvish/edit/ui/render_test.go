package ui

import "testing"

type testRenderer struct {
}

func (testRenderer) Render(b *Buffer) {
	b.WriteString("xy", "1")
}

func TestRender(t *testing.T) {
	b := Render(testRenderer{}, 10)
	if b.Width != 10 {
		t.Errorf("Rendered Buffer has Width %d, want %d", b.Width, 10)
	}
	if eq, _ := CompareCells(b.Lines[0], []Cell{{"x", 1, "1"}, {"y", 1, "1"}}); !eq {
		t.Errorf("Rendered Buffer has unexpected content")
	}
}
