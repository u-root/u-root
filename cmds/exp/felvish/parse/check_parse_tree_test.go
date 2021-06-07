package parse

import "fmt"

// checkParseTree checks whether the parse tree part of a Node is well-formed.
func checkParseTree(n Node) error {
	children := n.Children()
	if len(children) == 0 {
		return nil
	}

	// Parent pointers of all children should point to me.
	for i, ch := range children {
		if ch.Parent() != n {
			return fmt.Errorf("parent of child %d (%s) is wrong: %s", i, summary(ch), summary(n))
		}
	}

	// The Begin of the first child should be equal to mine.
	if children[0].Begin() != n.Begin() {
		return fmt.Errorf("gap between node and first child: %s", summary(n))
	}
	// The End of the last child should be equal to mine.
	nch := len(children)
	if children[nch-1].End() != n.End() {
		return fmt.Errorf("gap between node and last child: %s", summary(n))
	}
	// Consecutive children have consecutive position ranges.
	for i := 0; i < nch-1; i++ {
		if children[i].End() != children[i+1].Begin() {
			return fmt.Errorf("gap between child %d and %d of: %s", i, i+1, summary(n))
		}
	}

	// Check children recursively.
	for _, ch := range n.Children() {
		err := checkParseTree(ch)
		if err != nil {
			return err
		}
	}
	return nil
}
